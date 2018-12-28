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

	ic, err := identicon.New("oh-hello@my-identicon.com", nil)
	if err != nil {
		panic(err.Error())
	}

	png.Encode(fi, ic.GenerateImage())
}
