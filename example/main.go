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

	ic, _ := identicon.New("identicon.png")

	// fmt.Println(hex.EncodeToString(ic.Hash))
	png.Encode(fi, *ic.GenerateImage())
}
