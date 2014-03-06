package edid

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

var TimingBitMap1 = []string{
	"720×400 @ 70 Hz",
	"720×400 @ 88 Hz",
	"640×480 @ 60 Hz",
	"640×480 @ 67 Hz",
	"640×480 @ 72 Hz",
	"640×480 @ 75 Hz",
	"800×600 @ 56 Hz",
	"800×600 @ 60 Hz"}

var TimingBitMap2 = []string{
	"800×600 @ 72 Hz",
	"800×600 @ 75 Hz",
	"832×624 @ 75 Hz",
	"1024×768 @ 87 Hz, interlaced (1024×768i)",
	"1024×768 @ 60 Hz",
	"1024×768 @ 72 Hz",
	"1024×768 @ 75 Hz",
	"1280×1024 @ 75 Hz"}

type MonitorRangeLimitDescriptor struct {
	MinimumVerticalFieldRate  byte
	MaximumVerticalFieldRate  byte
	MinimumHorizontalLineRate byte
	MaximumHorizontalLineRate byte
	MaximumPixelClockRate     byte
}

type DetailedTimingDescriptor struct {
	PixelClock               uint32
	HorizontalActive         uint16
	HorizontalBlanking       uint16
	VerticalActive           uint16
	VerticalBlanking         uint16
	HorizontalSyncOffset     uint16
	HorizontalSyncPulseWidth uint16
	VerticalSyncOffset       uint16
	VerticalSyncPulseWidth   uint16
	HorizontalImageSize      uint16
	VerticalImageSize        uint16
	HorizontalBorder         byte
	VerticalBorder           byte
	Interlaced               bool
	Stereo                   string
	SyncType                 byte
	HorizontalSyncPolarity   bool
	VerticalSyncPolarity     bool
}

type Edid struct {
	CheckSum          bool
	Header            [8]byte
	ManufacturerId    string
	ProductCode       uint16
	SerialNumber      uint32
	WeekOfManufacture byte
	YearOfManufacture uint
	EdidVersion       byte
	EdidRevision      byte

	DigitalInput      bool
	VESADFPCompatible bool

	CompositeSyncSupported bool
	SyncOnGreenSupported   bool

	MaximumHorizontalImageSize byte
	MaximumVerticalImageSize   byte
	DisplayGamma               float32

	DPMSStandbySupported   bool
	DPMSSuspendSupported   bool
	DPMSActiveOffSupported bool
	DisplayType            string

	RedX   float64
	RedY   float64
	GreenX float64
	GreenY float64
	BlueX  float64
	BlueY  float64
	WhiteX float64
	WhiteY float64

	TimingBitMap1 byte
	TimingBitMap2 byte
	TimingBitMap3 byte

	StandardTimingInformation []string

	DetailedTimingDescriptors []DetailedTimingDescriptor

	MonitorName         string
	MonitorSerialNumber string

	MonitorRangeLimitDescriptors []MonitorRangeLimitDescriptor

	NumberOfExtensions byte
}

func NewEdid(edidBytes []byte) (*Edid, error) {
	edid := new(Edid)
	for i := 0; i < 8; i++ {
		edid.Header[i] = edidBytes[i]
	}

	//check the checksum
	var checkSum byte
	checkSum = 0
	for i := 0; i < 128; i++ {
		checkSum += edidBytes[i]
	}
	if checkSum == 0 {
		edid.CheckSum = true
	} else {
		edid.CheckSum = false
	}

	manufacturerId := edidBytes[8:10]
	edid.ManufacturerId = fmt.Sprintf("%c%c%c", (manufacturerId[0]>>2&0x1f)+'A'-1, (((manufacturerId[0]&0x3)<<3)|((manufacturerId[1]&0xe0)>>5))+'A'-1, (manufacturerId[1]&0x1f)+'A'-1)

	p := bytes.NewBuffer(edidBytes)
	//skip header and manufacturer id
	p.Next(10)

	var productCode uint16
	binary.Read(p, binary.LittleEndian, &productCode)
	edid.ProductCode = productCode

	var serialNumber uint32
	binary.Read(p, binary.LittleEndian, &serialNumber)
	edid.SerialNumber = serialNumber

	var weekOfManufacture byte
	binary.Read(p, binary.LittleEndian, &weekOfManufacture)
	edid.WeekOfManufacture = weekOfManufacture

	var yearOfManufacture byte
	binary.Read(p, binary.LittleEndian, &yearOfManufacture)
	edid.YearOfManufacture = uint(yearOfManufacture) + 1990

	var edidVersion byte
	binary.Read(p, binary.LittleEndian, &edidVersion)
	edid.EdidVersion = edidVersion

	var edidRevision byte
	binary.Read(p, binary.LittleEndian, &edidRevision)
	edid.EdidRevision = edidRevision

	var displayParam byte
	binary.Read(p, binary.LittleEndian, &displayParam)
	if displayParam&0x80 > 0 {
		edid.DigitalInput = true
		if displayParam&0x1 > 0 {
			edid.VESADFPCompatible = true
		}
	} else {
		edid.DigitalInput = false
		if displayParam&0x4 > 0 {
			edid.CompositeSyncSupported = true
		}
		if displayParam&0x2 > 0 {
			edid.SyncOnGreenSupported = true
		}

	}

	var horizontalDisplaySize byte
	binary.Read(p, binary.LittleEndian, &horizontalDisplaySize)
	edid.MaximumHorizontalImageSize = horizontalDisplaySize

	var verticalDisplaySize byte
	binary.Read(p, binary.LittleEndian, &verticalDisplaySize)
	edid.MaximumVerticalImageSize = verticalDisplaySize

	var displayGamma byte
	binary.Read(p, binary.LittleEndian, &displayGamma)
	edid.DisplayGamma = (float32(displayGamma) / 100) + 1

	var featureBitmap byte
	binary.Read(p, binary.LittleEndian, &featureBitmap)
	if featureBitmap&0x80 > 0 {
		edid.DPMSStandbySupported = true
	}
	if featureBitmap&0x40 > 0 {
		edid.DPMSSuspendSupported = true
	}
	if featureBitmap&0x20 > 0 {
		edid.DPMSActiveOffSupported = true
	}

	dispType := (featureBitmap & 0x18) >> 3

	if edid.DigitalInput {
		switch dispType {
		case 0:
			edid.DisplayType = "RGB 4:4:4"
		case 1:
			edid.DisplayType = "RGB 4:4:4 + YCrCb 4:4:4"
		case 2:
			edid.DisplayType = "RGB 4:4:4 + YCrCb 4:2:2"
		case 3:
			edid.DisplayType = "RGB 4:4:4 + YCrCb 4:4:4 + YCrCb 4:2:2"
		}
	} else {
		switch dispType {
		case 0:
			edid.DisplayType = "Monochrome/Grayscale"
		case 1:
			edid.DisplayType = "RGB color"
		case 2:
			edid.DisplayType = "Non-RGB color"
		case 3:
			edid.DisplayType = "Undefined"
		}
	}

	//25
	var redGreenLSB byte
	binary.Read(p, binary.LittleEndian, &redGreenLSB)

	//26
	var blueWhiteLSB byte
	binary.Read(p, binary.LittleEndian, &blueWhiteLSB)

	//27
	var redXMSB byte
	binary.Read(p, binary.LittleEndian, &redXMSB)

	//28
	var redYMSB byte
	binary.Read(p, binary.LittleEndian, &redYMSB)

	//29
	var greenXMSB byte
	binary.Read(p, binary.LittleEndian, &greenXMSB)

	//30
	var greenYMSB byte
	binary.Read(p, binary.LittleEndian, &greenYMSB)

	//31
	var blueXMSB byte
	binary.Read(p, binary.LittleEndian, &blueXMSB)

	//32
	var blueYMSB byte
	binary.Read(p, binary.LittleEndian, &blueYMSB)

	//33
	var whiteXMSB byte
	binary.Read(p, binary.LittleEndian, &whiteXMSB)

	//34
	var whiteYMSB byte
	binary.Read(p, binary.LittleEndian, &whiteYMSB)

	edid.RedX = float64((uint16(redXMSB)<<2)|((uint16(redGreenLSB)>>6)&0x3)) / 1024
	edid.RedY = float64((uint16(redYMSB)<<2)|((uint16(redGreenLSB)>>4)&0x3)) / 1024
	edid.GreenX = float64((uint16(greenXMSB)<<2)|((uint16(redGreenLSB)>>2)&0x3)) / 1024
	edid.GreenY = float64((uint16(greenYMSB)<<2)|(uint16(redGreenLSB)&0x3)) / 1024
	edid.BlueX = float64((uint16(blueXMSB)<<2)|((uint16(blueWhiteLSB)>>6)&0x3)) / 1024
	edid.BlueY = float64((uint16(blueYMSB)<<2)|((uint16(blueWhiteLSB)>>4)&0x3)) / 1024
	edid.WhiteX = float64((uint16(whiteXMSB)<<2)|((uint16(blueWhiteLSB)>>2)&0x3)) / 1024
	edid.WhiteY = float64((uint16(whiteYMSB)<<2)|(uint16(blueWhiteLSB)&0x3)) / 1024

	binary.Read(p, binary.LittleEndian, &edid.TimingBitMap1)
	binary.Read(p, binary.LittleEndian, &edid.TimingBitMap2)
	binary.Read(p, binary.LittleEndian, &edid.TimingBitMap3)

	for i := 0; i < 8; i++ {
		var temp byte
		binary.Read(p, binary.LittleEndian, &temp)
		xResolution := (uint(temp) + 31) * 8
		binary.Read(p, binary.LittleEndian, &temp)
		pixelRatio := temp >> 6
		verticalFrequency := (uint(temp) & 63) + 60
		var yResolution uint
		var pixelRatioStr string

		switch pixelRatio {
		case 0:
			//16:10
			yResolution = xResolution * 10 / 16
			pixelRatioStr = "16:10"
		case 1:
			//4:3
			pixelRatioStr = "4:3"
			yResolution = xResolution * 3 / 4
		case 2:
			//5:4
			pixelRatioStr = "5:4"
			yResolution = xResolution * 4 / 5
		case 3:
			//16:9
			pixelRatioStr = "16:9"
			yResolution = xResolution * 9 / 16

		}
		edid.StandardTimingInformation = append(edid.StandardTimingInformation, fmt.Sprintf("%dx%d %s @%dHz", xResolution, yResolution, pixelRatioStr, verticalFrequency))

	}

	//read 4 descriptors

	for i := 0; i < 4; i++ {
		var temp [18]byte
		var descriptor DetailedTimingDescriptor
		binary.Read(p, binary.LittleEndian, &temp)
		descriptor.PixelClock = ((uint32(temp[1]) << 8) | uint32(temp[0])) * 10
		descriptor.HorizontalActive = (uint16(temp[4]&240) << 4) | uint16(temp[2])
		descriptor.HorizontalBlanking = (uint16(temp[4]&15) << 8) | uint16(temp[3])
		descriptor.VerticalActive = (uint16(temp[7]&240) << 4) | uint16(temp[5])
		descriptor.VerticalBlanking = (uint16(temp[7]&15) << 8) | uint16(temp[6])
		descriptor.HorizontalSyncOffset = ((uint16(temp[11]) & 192) << 2) | uint16(temp[8])
		descriptor.HorizontalSyncPulseWidth = ((uint16(temp[11]) & 48) << 4) | uint16(temp[9])
		descriptor.VerticalSyncOffset = ((uint16(temp[11]) & 12) << 2) | ((uint16(temp[10]) & 240) >> 4)
		descriptor.VerticalSyncPulseWidth = ((uint16(temp[11]) & 3) << 4) | (uint16(temp[10]) & 15)
		descriptor.HorizontalImageSize = ((uint16(temp[14]) & 240) << 4) | uint16(temp[12])
		descriptor.VerticalImageSize = ((uint16(temp[14]) & 15) << 8) | uint16(temp[13])
		descriptor.HorizontalBorder = temp[15]
		descriptor.VerticalBorder = temp[16]
		descriptor.Interlaced = (temp[17] & 128) > 0
		zeroBit := temp[17] & 1
		stereoMode := (temp[17] & 96) >> 5
		if stereoMode == 0 {
			descriptor.Stereo = "No Stereo"
		} else {
			if zeroBit == 1 {
				switch stereoMode {
				case 1:
					descriptor.Stereo = "2-way interleaved stereo - Right image on even lines"
				case 2:
					descriptor.Stereo = "2-way interleaved stereo - Left image on even lines"
				case 3:
					descriptor.Stereo = "2-way interleaved stereo - side-by-side"
				}

			} else {
				switch stereoMode {
				case 1:
					descriptor.Stereo = "Field sequential, sync=1 during right"
				case 2:
					descriptor.Stereo = "similar, sync=1 during left"
				case 3:
					descriptor.Stereo = "4-way interleaved stereo"
				}
			}
		}

		descriptor.SyncType = (temp[17] & 24) >> 3
		descriptor.VerticalSyncPolarity = ((temp[17] & 4) >> 2) > 0
		descriptor.HorizontalSyncPolarity = ((temp[17] & 2) >> 1) > 0

		if descriptor.PixelClock == 0 {
			descriptorType := temp[3]
			switch descriptorType {
			case 0xFF:
				edid.MonitorSerialNumber = string(temp[5:])
			case 0xFC:
				mon := string(temp[5:])
				mon = strings.Replace(mon, "\n", "", -1)
				mon = strings.TrimSpace(mon)
				edid.MonitorName = mon
			case 0xFD:
				//Monitor range limits
				var monitorRangeDescriptor MonitorRangeLimitDescriptor
				monitorRangeDescriptor.MinimumVerticalFieldRate = temp[5]
				monitorRangeDescriptor.MaximumVerticalFieldRate = temp[6]
				monitorRangeDescriptor.MinimumHorizontalLineRate = temp[7]
				monitorRangeDescriptor.MaximumHorizontalLineRate = temp[8]
				monitorRangeDescriptor.MaximumPixelClockRate = temp[9] * 10
				edid.MonitorRangeLimitDescriptors = append(edid.MonitorRangeLimitDescriptors, monitorRangeDescriptor)
			}
		} else {
			edid.DetailedTimingDescriptors = append(edid.DetailedTimingDescriptors, descriptor)
		}
	}

	binary.Read(p, binary.LittleEndian, &edid.NumberOfExtensions)
	return edid, nil
}

func (edid *Edid) PrintableHeader() string {
	var buffer bytes.Buffer
	buffer.WriteString("0x")
	for i := 0; i < len(edid.Header); i++ {
		s := fmt.Sprintf("%02X", edid.Header[i])
		buffer.WriteString(s)
	}
	return buffer.String()
}

func (e *Edid) PrettyPrint() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 30, 20, 4, '-', 0)
	fmt.Fprintln(w, "Valid Checksum: \t ", e.CheckSum)
	fmt.Fprintln(w, "Header: \t ", e.PrintableHeader())
	fmt.Fprintln(w, "Monitor Name: \t ", e.MonitorName)
	fmt.Fprintln(w, "Monitor Serial Number: \t ", e.MonitorSerialNumber)
	fmt.Fprintln(w, "Manufacturer Name: \t ", e.ManufacturerId)
	fmt.Fprintln(w, "Product Code: \t ", e.ProductCode)
	fmt.Fprintln(w, "Serial Number: \t ", e.SerialNumber)
	fmt.Fprintln(w, "Week of Manufacture: \t ", e.WeekOfManufacture)
	fmt.Fprintln(w, "Year of Manufacture: \t ", e.YearOfManufacture)
	fmt.Fprintln(w, "EDID Version: \t ", e.EdidVersion)
	fmt.Fprintln(w, "EDID Revision: \t ", e.EdidRevision)
	fmt.Fprintln(w, "\nBasic display parameters:\n")

	if e.DigitalInput {
		fmt.Fprintln(w, "    Video Input Definition: \t ", "Digital")
		fmt.Fprintln(w, "    VESA DFP Compatibility: \t ", e.VESADFPCompatible)
	} else {
		fmt.Fprintln(w, "    Video Input Definition: \t ", "Analog")
		fmt.Fprintln(w, "    Composite Sync Supported: \t ", e.CompositeSyncSupported)
		fmt.Fprintln(w, "    Sync on Green Supported: \t ", e.SyncOnGreenSupported)
	}
	fmt.Fprintln(w, "    Max Horizontal Image Size: \t ", e.MaximumHorizontalImageSize, "cm")
	fmt.Fprintln(w, "    Max Vertical Image Size: \t ", e.MaximumVerticalImageSize, "cm")
	fmt.Fprintln(w, "    Display Gamma: \t ", e.DisplayGamma)

	fmt.Fprintln(w, "\nPower Management:")
	fmt.Fprintln(w, "    DPMS Standby Supported: \t ", e.DPMSStandbySupported)
	fmt.Fprintln(w, "    DPMS Suspend Supported: \t ", e.DPMSSuspendSupported)
	fmt.Fprintln(w, "    DPMS Active Off Supported: \t ", e.DPMSActiveOffSupported)
	fmt.Fprintln(w, "    Display Type: \t ", e.DisplayType)

	fmt.Fprintln(w, "\nChroma Information:")
	fmt.Fprintln(w, "    Red X: \t ", strconv.FormatFloat(e.RedX, 'f', 3, 64))
	fmt.Fprintln(w, "    Red Y: \t ", strconv.FormatFloat(e.RedY, 'f', 3, 64))
	fmt.Fprintln(w, "    Green X: \t ", strconv.FormatFloat(e.GreenX, 'f', 3, 64))
	fmt.Fprintln(w, "    Green Y: \t ", strconv.FormatFloat(e.GreenY, 'f', 3, 64))
	fmt.Fprintln(w, "    Blue X: \t ", strconv.FormatFloat(e.BlueX, 'f', 3, 64))
	fmt.Fprintln(w, "    Blue Y: \t ", strconv.FormatFloat(e.BlueY, 'f', 3, 64))
	fmt.Fprintln(w, "    White X: \t ", strconv.FormatFloat(e.WhiteX, 'f', 3, 64))
	fmt.Fprintln(w, "    White Y: \t ", strconv.FormatFloat(e.WhiteY, 'f', 3, 64))

	fmt.Fprintln(w, "\nTimings Bitmaps:")

	counter := 0
	var mask byte
	for mask = 0x80; mask != 0; mask >>= 1 {
		if e.TimingBitMap1&mask > 0 {
			fmt.Fprintln(w, "    "+TimingBitMap1[counter])
		}
		counter++
	}
	counter = 0
	for mask = 0x80; mask != 0; mask >>= 1 {
		if e.TimingBitMap2&mask > 0 {
			fmt.Fprintln(w, "    "+TimingBitMap2[counter])
		}
		counter++
	}

	if e.TimingBitMap3&0x80 > 0 {
		fmt.Fprintln(w, "    1152x870 @ 75 Hz (Apple Macintosh II)")
	}
	fmt.Fprintln(w, "\nStandard Timing Identification:")
	for _, element := range e.StandardTimingInformation {
		fmt.Fprintln(w, "    ", element)
	}

	for index, element := range e.DetailedTimingDescriptors {
		fmt.Fprintln(w, "\nDetailed Timing/Descriptor block", (index + 1))
		fmt.Fprintln(w, "    Pixel Clock: \t ", element.PixelClock, "kHz")
		fmt.Fprintln(w, "    Horizontal Active: \t ", element.HorizontalActive, "pixels")
		fmt.Fprintln(w, "    Horizontal Blanking: \t ", element.HorizontalBlanking, "pixels")
		fmt.Fprintln(w, "    Vertical Active: \t ", element.VerticalActive, "pixels")
		fmt.Fprintln(w, "    Vertical Blanking: \t ", element.VerticalBlanking, "pixels")
		fmt.Fprintln(w, "    Horizontal Sync Offset: \t ", element.HorizontalSyncOffset, "pixels")
		fmt.Fprintln(w, "    Horizontal Sync Pulse Width: \t ", element.HorizontalSyncPulseWidth, "pixels")
		fmt.Fprintln(w, "    Vertical Sync Offset: \t ", element.VerticalSyncOffset, "lines")
		fmt.Fprintln(w, "    Vertical Sync Pulse Width: \t ", element.VerticalSyncPulseWidth, "lines")
		fmt.Fprintln(w, "    Horizontal Image Size: \t ", element.HorizontalImageSize, "mm")
		fmt.Fprintln(w, "    Vertical Image Size: \t ", element.VerticalImageSize, "mm")
		fmt.Fprintln(w, "    Horizontal Border: \t ", element.HorizontalBorder, "px each side")
		fmt.Fprintln(w, "    Vertical Border: \t ", element.VerticalBorder, "px each side")
		fmt.Fprintln(w, "    Interlaced: \t ", element.Interlaced)
		fmt.Fprintln(w, "    Stereo Mode: \t ", element.Stereo)
		switch element.SyncType {
		case 0:
			fmt.Fprintln(w, "    Sync Type: \t ", "Analog composite")
		case 1:
			fmt.Fprintln(w, "    Sync Type: \t ", "Bipolar analog composite")
		case 2:
			fmt.Fprintln(w, "    Sync Type: \t ", "Digital composite (on HSync)")
		case 3:
			fmt.Fprintln(w, "    Sync Type: \t ", "Digital separate")
			fmt.Fprintln(w, "    Vertical Sync Polarity: \t ", element.VerticalSyncPolarity)
			fmt.Fprintln(w, "    Horizontal Sync Polarity: \t ", element.HorizontalSyncPolarity)
		}

	}

	for index, element := range e.MonitorRangeLimitDescriptors {
		fmt.Fprintln(w, "\nMonitor range limits descriptor block", (index + 1))
		fmt.Fprintln(w, "    Minimum Vertical Field Rate: \t ", element.MinimumVerticalFieldRate, "Hz")
		fmt.Fprintln(w, "    Maximum Vertical Field Rate: \t ", element.MaximumVerticalFieldRate, "Hz")
		fmt.Fprintln(w, "    Minimum Horizontal Line Rate: \t ", element.MinimumHorizontalLineRate, "kHz")
		fmt.Fprintln(w, "    Maximum Horizontal Line Rate: \t ", element.MaximumHorizontalLineRate, "kHz")
		fmt.Fprintln(w, "    Maximum Pixel Clock Rate: \t ", element.MaximumPixelClockRate, "MHz")
	}
	fmt.Fprintln(w, "\nTotal Number of Extensions: \t ", e.NumberOfExtensions, "(not parsed)")

	fmt.Fprintln(w, "\n**************************************************************************")
	fmt.Fprintln(w, "*                Bugs? Contact anoopengineer@gmail.com                   *")
	fmt.Fprintln(w, "**************************************************************************\n")
	w.Flush()
}
