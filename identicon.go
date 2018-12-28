package identicon

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
)

const (
	tilesPerDimension = 5
)

var (
	defaultBackgroundColor = RGB(240, 240, 240)
	defaultImageSize       = 100
)

// Identicon defines an identicon
type Identicon struct {
	ID      string
	Hash    []byte
	Tiles   [][]bool
	Color   color.Color
	Options *Options
}

// Options control some inner mechanics
type Options struct {
	BackgroundColor color.NRGBA
	Debug           bool
	ImageSize       int
}

// New returns a new identicon based on given ID string
func New(ID string, opts *Options) (*Identicon, error) {

	if opts == nil {
		opts = &Options{
			BackgroundColor: defaultBackgroundColor,
		}
	}

	if opts.ImageSize <= 0 {
		opts.ImageSize = defaultImageSize
	}

	return &Identicon{
		ID:      ID,
		Hash:    MD5(ID),
		Options: opts,
	}, nil
}

// GenerateImage returns an generated Image representation of the identicon
func (ic *Identicon) GenerateImage() *image.RGBA {

	ic.populateTiles()
	ic.defineColor()
	if ic.Options.Debug {
		ic.debugPrintTiles()
	}

	bounds := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{ic.Options.ImageSize, ic.Options.ImageSize},
	}

	img := image.NewRGBA(bounds)

	// Background fill
	draw.Draw(img, img.Bounds(), &image.Uniform{ic.Options.BackgroundColor}, image.ZP, draw.Src)

	// Iterate tiles and draw
	for xTile := 0; xTile < tilesPerDimension; xTile++ {
		for yTile := 0; yTile < tilesPerDimension; yTile++ {
			if ic.Tiles[xTile][yTile] {
				ic.drawTile(img, xTile, yTile)
			}

		}
	}

	return img
}

func (ic *Identicon) drawTile(img *image.RGBA, xTile, yTile int) {

	xStart := (xTile * (ic.Options.ImageSize / tilesPerDimension))
	xEnd := xStart + (ic.Options.ImageSize / tilesPerDimension)

	yStart := (yTile * (ic.Options.ImageSize / tilesPerDimension))
	yEnd := yStart + (ic.Options.ImageSize / tilesPerDimension)

	// fmt.Println("x", xStart, xEnd)
	// fmt.Println("y", yStart, yEnd)

	bounds := image.Rect(xStart, yStart, xEnd, yEnd)
	draw.Draw(img, bounds, &image.Uniform{ic.Color}, image.ZP, draw.Src)
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

			if ic.Options.Debug {
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

// debugPrintTiles prints the tiles at positions x,y
func (ic *Identicon) debugPrintTiles() {
	for x := range ic.Tiles {
		for y, v := range ic.Tiles[x] {
			fmt.Printf("Tile %d:%d = %v\n", x, y, v)
		}
	}
}

// MD5 returns MD5 hash of given input string as byte slice
func MD5(text string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hasher.Sum(nil)
}

// RGB returns color.NRGBA struct for given red, green and blue values
func RGB(r, g, b uint8) color.NRGBA {
	return color.NRGBA{R: r, G: g, B: b, A: 255}
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
