package main

import (
	"encoding/hex"
	"fmt"
	"github.com/anoopengineer/edidparser/edid"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please pass the input filename as an argument")
		os.Exit(1)
	}
	fileName := os.Args[1]
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Unable to read file", fileName)
		os.Exit(1)
	}

	str := string(bs)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\r\n", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.TrimSpace(str)
	//fmt.Printf("%s\n", str)

	edidBytes, _ := hex.DecodeString(str)
	printEdidBytes(edidBytes)
	e, err := edid.NewEdid(edidBytes)
	//fmt.Println("Man - " + e.PrintableManufacturerId())

	fmt.Println()
	if err != nil {
		log.Fatal("Unable to parse EDID ", err)
	} else {
		e.PrettyPrint()
	}

}

func printEdidBytes(edid []byte) {
	fmt.Println("EDID dump")
	columnCounter := 0
	for i := 0; i < len(edid); i++ {
		fmt.Printf("%02X ", edid[i])
		if columnCounter == 15 {
			fmt.Println()
			columnCounter = 0
		} else {
			columnCounter++
		}
	}
}
