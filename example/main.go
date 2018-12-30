package main

import (
	"image/png"
	"log"
	"os"

	"github.com/tsdtsdtsd/identicon"
)

func main() {

	fi, err := os.Create("identicon.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	ic, err := identicon.New("oh-hello@my-identicon.com", &identicon.Options{Debug: false, BackgroundColor: identicon.RGB(240, 240, 240)})
	if err != nil {
		panic(err.Error())
	}

	png.Encode(fi, ic.GenerateImage())
}
