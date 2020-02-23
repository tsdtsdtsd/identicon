package main

import (
	"image/png"
	"log"
	"os"

	"github.com/tsdtsdtsd/identicon"
)

func main() {

	// Just create a new identicon
	ic, err := identicon.New(

		// The identicon ID string is mandatory.
		// Same string will always result in the same generated identicon.
		// Typically this is a username or email address.
		"jack@example.com",
		// You can define custom options or pass nil for defaults
		&identicon.Options{
			BackgroundColor: identicon.RGB(240, 240, 240),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Now you are free to use identicon like any other image.Image or draw.Image interface
	fi, err := os.Create("identicon.png")
	if err != nil {
		log.Fatal(err)
	}

	png.Encode(fi, ic)
	fi.Close()
}
