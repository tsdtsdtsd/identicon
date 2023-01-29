package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/tsdtsdtsd/identicon/v1"
)

var genIDs = []string{
	"identicon",
	"test-string",
	"Amazatron3000",
	"yay-identicons",
	"m.jackson",
	"12monkeys",
	"Stan.Lee",
	"gogopher",
	"notblue",
	"test",
}

func main() {

	// Prepare file
	fi, err := os.Create("identicon-banner.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	var (
		spacing = 15
		size    = 100
		half    = len(genIDs) / 2
		x1      = (half * size) + ((half + 1) * spacing)
		y1      = (2 * size) + (3 * spacing)
	)

	// Create banner and fill background
	banner := image.NewRGBA(image.Rect(0, 0, x1, y1))
	draw.Draw(banner, image.Rect(0, 0, x1, y1), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	// Iterate IDs
	for i, id := range genIDs {

		// Create a new identicon
		ic, err := identicon.New(id)
		if err != nil {
			log.Fatal(err)
		}

		// Calculate position on banner
		x0 := ((i + 1) * spacing) + (i * size)
		y0 := spacing
		if i >= half {
			x0 = ((i - half + 1) * spacing) + ((i - half) * size)
			y0 += spacing + size
		}

		x1 := x0 + size
		y1 := y0 + size

		// Draw identicon on banner
		draw.Draw(banner, image.Rect(x0, y0, x1, y1), ic, image.Point{}, draw.Src)
	}

	// Encode and write banner to file
	png.Encode(fi, banner)
}
