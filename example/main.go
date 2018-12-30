package main

import (
	"image/png"
	"log"
	"os"

	"github.com/tsdtsdtsd/identicon"
)

func main() {

	fi, err := os.Create("test-string.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	ic, err := identicon.New("test-string", &identicon.Options{Debug: false, BackgroundColor: identicon.RGB(235, 235, 235)})
	if err != nil {
		panic(err.Error())
	}

	png.Encode(fi, ic.GenerateImage())
}
