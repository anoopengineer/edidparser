EDID parser in Go Language
---------------------------

This is an EDID (http://en.wikipedia.org/wiki/Extended_display_identification_data) parser library and application written in Go language. 

This contains an executable (currently available for windows - edidparser.exe and can be easily compiled for other platforms as well) and a library that can be imported and used to parse edids. 

Binary usage:
---------------------------

First dump the edid bytes to a file (space separated bytes - see edid_input.txt file to see a sample input file). And pass the file path as argument to the edidparser binary.

eg: `edidparser.exe <path_to_file_containing_edid_dump>`

	C:\>edidparser.exe edid_input.txt
	EDID dump
	00 FF FF FF FF FF FF 00 41 0C 00 00 00 00 00 00 
	00 15 01 03 80 00 00 78 0A 2F CD A3 54 45 97 24 
	11 46 47 21 08 00 01 01 01 01 01 01 01 01 01 01 
	01 01 01 01 01 01 02 3A 80 18 71 38 2D 40 58 2C 
	45 00 10 09 00 00 00 1E 8C 0A D0 8A 20 E0 2D 10 
	10 3E 96 00 10 09 00 00 00 18 00 00 00 FC 00 50 
	48 49 4C 49 50 53 0A 20 20 20 20 20 00 00 00 FD 
	00 17 3F 0F 45 0F 00 0A 20 20 20 20 20 20 01 68 
	02 03 20 70 4A 10 03 20 22 04 02 05 06 07 01 26 
	09 07 07 15 07 50 83 01 00 00 65 03 0C 00 10 00 
	01 1D 80 3E 73 38 2D 40 7E 2C 45 80 10 09 00 00 
	00 1E 01 1D 80 18 71 38 2D 40 58 2C 45 00 10 09 
	00 00 00 1E 01 1D 00 72 51 D0 1E 20 6E 28 55 00 
	10 09 00 00 00 1E 66 21 50 B0 51 00 1B 30 40 70 
	36 00 10 09 00 00 00 1E A9 1A 00 A0 50 00 16 30 
	30 20 37 00 05 03 00 00 00 1A 00 00 00 00 00 A9 

	Valid Checksum: --------------  true
	Header: ----------------------  0x00FFFFFFFFFFFF00
	Monitor Name: ----------------  PHILIPS
	Monitor Serial Number: -------  
	Manufacturer Name: -----------  PHL
	Product Code: ----------------  0
	Serial Number: ---------------  0
	Week of Manufacture: ---------  0
	Year of Manufacture: ---------  2011
	EDID Version: ----------------  1
	EDID Revision: ---------------  3

	Basic display parameters:

	    Video Input Definition: -------  Digital
	    VESA DFP Compatibility: -------  false
	    Max Horizontal Image Size: ----  0 cm
	    Max Vertical Image Size: ------  0 cm
	    Display Gamma: ----------------  2.2

	Power Management:
	    DPMS Standby Supported: -------  false
	    DPMS Suspend Supported: -------  false
	    DPMS Active Off Supported: ----  false
	    Display Type: -----------------  RGB 4:4:4 + YCrCb 4:4:4

	Chroma Information:
	    Red X: -------------------  0.637
	    Red Y: -------------------  0.330
	    Green X: -----------------  0.272
	    Green Y: -----------------  0.593
	    Blue X: ------------------  0.144
	    Blue Y: ------------------  0.066
	    White X: -----------------  0.276
	    White Y: -----------------  0.278

	Timings Bitmaps:
	    640×480 @ 60 Hz
	    800×600 @ 60 Hz
	    1024×768 @ 60 Hz

	Standard Timing Identification:
	     256x160 16:10 @61Hz
	     256x160 16:10 @61Hz
	     256x160 16:10 @61Hz
	     256x160 16:10 @61Hz
	     256x160 16:10 @61Hz
	     256x160 16:10 @61Hz
	     256x160 16:10 @61Hz
	     256x160 16:10 @61Hz

	Detailed Timing/Descriptor block 1
	    Pixel Clock: --------------------  148500 kHz
	    Horizontal Active: --------------  1920 pixels
	    Horizontal Blanking: ------------  280 pixels
	    Vertical Active: ----------------  1080 pixels
	    Vertical Blanking: --------------  45 pixels
	    Horizontal Sync Offset: ---------  88 pixels
	    Horizontal Sync Pulse Width: ----  44 pixels
	    Vertical Sync Offset: -----------  4 lines
	    Vertical Sync Pulse Width: ------  5 lines
	    Horizontal Image Size: ----------  16 mm
	    Vertical Image Size: ------------  9 mm
	    Horizontal Border: --------------  0 px each side
	    Vertical Border: ----------------  0 px each side
	    Interlaced: ---------------------  false
	    Stereo Mode: --------------------  No Stereo
	    Sync Type: ----------------------  Digital separate
	    Vertical Sync Polarity: ---------  true
	    Horizontal Sync Polarity: -------  true

	Detailed Timing/Descriptor block 2
	    Pixel Clock: --------------------  27000 kHz
	    Horizontal Active: --------------  720 pixels
	    Horizontal Blanking: ------------  138 pixels
	    Vertical Active: ----------------  480 pixels
	    Vertical Blanking: --------------  45 pixels
	    Horizontal Sync Offset: ---------  16 pixels
	    Horizontal Sync Pulse Width: ----  62 pixels
	    Vertical Sync Offset: -----------  9 lines
	    Vertical Sync Pulse Width: ------  6 lines
	    Horizontal Image Size: ----------  16 mm
	    Vertical Image Size: ------------  9 mm
	    Horizontal Border: --------------  0 px each side
	    Vertical Border: ----------------  0 px each side
	    Interlaced: ---------------------  false
	    Stereo Mode: --------------------  No Stereo
	    Sync Type: ----------------------  Digital separate
	    Vertical Sync Polarity: ---------  false
	    Horizontal Sync Polarity: -------  false

	Monitor range limits descriptor block 1
	    Minimum Vertical Field Rate: -----  23 Hz
	    Maximum Vertical Field Rate: -----  63 Hz
	    Minimum Horizontal Line Rate: ----  15 kHz
	    Maximum Horizontal Line Rate: ----  69 kHz
	    Maximum Pixel Clock Rate: --------  150 MHz

	Total Number of Extensions: ----  1 (not parsed)

	**************************************************************************
	*                Bugs? Contact anoopengineer@gmail.com                   *
	**************************************************************************

