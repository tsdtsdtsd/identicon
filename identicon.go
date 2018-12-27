package identicon

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
)

const (
	imgWidth  = 100
	imgHeight = 100

	tilesPerDimension = 5
)

var (
	backgroundColor = rgb(235, 235, 235)

	debug = false
)

// Identicon defines an identicon
type Identicon struct {
	ID    string
	Hash  []byte
	Tiles [][]bool
	Color color.Color
}

// New returns a new identicon based on given ID string
func New(ID string) (*Identicon, error) {

	// A valid hash is mandatory
	hash, err := MD5(ID)
	if err != nil {
		return nil, err
	}

	return &Identicon{
		ID:   ID,
		Hash: hash,
	}, nil
}

// GenerateImage returns an generated Image representation of the identicon
func (ic *Identicon) GenerateImage() image.Image {

	ic.populateTiles()
	ic.defineColor()
	if debug {
		ic.DebugPrintTiles()
	}

	bounds := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{imgWidth, imgHeight},
	}

	img := image.NewPaletted(
		bounds,
		color.Palette{
			backgroundColor,
			ic.Color,
		},
	)

	// Iterate tiles and draw
	for xTile := 0; xTile < tilesPerDimension; xTile++ {
		for yTile := 0; yTile < tilesPerDimension; yTile++ {
			if ic.Tiles[xTile][yTile] {
				ic.drawTile(img, xTile, yTile)
			}

		}
	}

	return img.SubImage(bounds)
}

func (ic *Identicon) drawTile(img *image.Paletted, xTile, yTile int) {

	xStart := (xTile * (imgWidth / tilesPerDimension))
	if xStart < 0 {
		xStart = 0
	}
	xEnd := xStart + (imgWidth / tilesPerDimension) - 1

	yStart := (yTile * (imgHeight / tilesPerDimension))
	if yStart < 0 {
		yStart = 0
	}
	yEnd := yStart + (imgHeight / tilesPerDimension) - 1

	// fmt.Println("x", xStart, xEnd)
	// fmt.Println("y", yStart, yEnd)

	for x := xStart; x <= xEnd; x++ {
		for y := yStart; y <= yEnd; y++ {
			img.SetColorIndex(x, y, 1)
		}
	}
}

func (ic *Identicon) populateTiles() {

	tiles := make([][]bool, tilesPerDimension)
	for i := range tiles {
		tiles[i] = make([]bool, tilesPerDimension)
	}

	ic.Tiles = tiles

	// Per image, we have 5x5 tiles available.
	// First 15 bytes of hash define tiles:
	//   - first 10 are the two leftmost cols and get mirrored to the rightmost cols
	//   - next 5 for the middle col
	// Last byte for the pixel color

	// Left
	var i int8
	for i = 0; i < 10; i++ {
		ic.calcTile(i, ic.Hash[i])
	}

	// Middle
	for i = 10; i < 15; i++ {
		ic.calcTile(i, ic.Hash[i])
	}

	// Mirror to right
	ic.mirror()
}

func (ic *Identicon) calcTile(pos int8, b byte) {

	lever := (int(b) & 2) > 0
	x, y := posToXY(pos)

	ic.Tiles[x][y] = lever
}

func (ic *Identicon) mirror() {

	for x := 0; x <= 1; x++ {

		xi := tilesPerDimension - 1 - x // mirror offset for cols

		for y := 0; y < tilesPerDimension; y++ {

			ic.Tiles[xi][y] = ic.Tiles[x][y]

			if debug {
				fmt.Printf("Mirroring %d:%d to %d:%d (%v)\n", x, y, xi, y, ic.Tiles[x][y])
			}
		}
	}
}

func (ic *Identicon) defineColor() {

	colorIdx := int(ic.Hash[15])

	// @todo need a custom palette

	// Index out of range?
	for colorIdx > len(palette.WebSafe)-1 {
		colorIdx = colorIdx / 9 * 7
	}

	ic.Color = palette.WebSafe[colorIdx]
}

// HashString returns hash value as string
func (ic *Identicon) HashString() string {
	return hex.EncodeToString(ic.Hash)
}

// DebugPrintTiles prints the tiles at positions x,y
func (ic *Identicon) DebugPrintTiles() {
	for x := range ic.Tiles {
		for y, v := range ic.Tiles[x] {
			fmt.Printf("Tile %d:%d = %v\n", x, y, v)
		}
	}
}

// MD5 returns MD5 hash of given input string as byte slice
func MD5(text string) ([]byte, error) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(text))
	if err != nil {
		return []byte{}, err
	}

	return hasher.Sum(nil), nil
}

func posToXY(pos int8) (x, y int) {

	// The two leftmost cols
	if pos < 10 {
		if pos%2 != 0 {
			x = 1
		}
		y = int(float32(pos) / 2.0)
	} else {
		// Middle col
		x = 2
		y = int(float32(pos) / 3.0)
	}

	return
}

func rgb(r, g, b uint8) color.NRGBA {
	return color.NRGBA{R: r, G: g, B: b, A: 255}
}
