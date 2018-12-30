// Package identicon helps generating deterministic user identicons
package identicon

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

const (
	tilesPerDimension = 5
)

var (
	defaultBackgroundColor = RGB(240, 240, 240)
	defaultImageSize       = 100
)

// Options control some inner mechanics
type Options struct {
	BackgroundColor color.NRGBA
	Debug           bool
	ImageSize       int
}

// Identicon defines an identicon
type Identicon struct {
	Color   color.Color
	Hash    []byte
	ID      string
	Options *Options
	Tiles   [][]bool
	// Pix holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
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

	rect := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{opts.ImageSize, opts.ImageSize},
	}

	w, h := rect.Dx(), rect.Dy()
	buf := make([]uint8, 4*w*h)

	ic := &Identicon{
		ID:      ID,
		Hash:    MD5(ID),
		Options: opts,
		Pix:     buf,
		Rect:    rect,
		Stride:  4 * w,
	}

	ic.generateImage()
	return ic, nil
}

// ColorModel returns the Image's color model.
func (ic *Identicon) ColorModel() color.Model {
	return color.ModelFunc(
		func(c color.Color) color.Color {
			// @todo
			return c
		},
	)
}

// Bounds returns the domain for which At can return non-zero color.
func (ic *Identicon) Bounds() image.Rectangle {
	return ic.Rect
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (ic *Identicon) At(x, y int) color.Color {
	return ic.NRGBAAt(x, y)
}

// NRGBAAt returns the color of the pixel at (x, y) as color.NRGBA.
func (ic *Identicon) NRGBAAt(x, y int) color.NRGBA {
	if !(image.Point{x, y}.In(ic.Rect)) {
		return color.NRGBA{}
	}
	i := ic.PixOffset(x, y)
	return color.NRGBA{ic.Pix[i+0], ic.Pix[i+1], ic.Pix[i+2], ic.Pix[i+3]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (ic *Identicon) PixOffset(x, y int) int {
	return (y-ic.Rect.Min.Y)*ic.Stride + (x-ic.Rect.Min.X)*4
}

// Set stores given color at position (x, y).
func (ic *Identicon) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(ic.Rect)) {
		return
	}

	i := ic.PixOffset(x, y)
	c1 := ic.ColorModel().Convert(c).(color.NRGBA)
	ic.Pix[i+0] = c1.R
	ic.Pix[i+1] = c1.G
	ic.Pix[i+2] = c1.B
	ic.Pix[i+3] = c1.A
}

// HashString returns hash value as string
func (ic *Identicon) HashString() string {
	return hex.EncodeToString(ic.Hash)
}

// generateImage generates image.Image representation of the identicon
func (ic *Identicon) generateImage() {

	ic.populateTiles()
	ic.defineColor()
	if ic.Options.Debug {
		ic.debugPrintTiles()
	}

	// Background fill
	draw.Draw(ic, ic.Bounds(), &image.Uniform{ic.Options.BackgroundColor}, image.ZP, draw.Src)

	// Iterate tiles and draw
	for xTile := 0; xTile < tilesPerDimension; xTile++ {
		for yTile := 0; yTile < tilesPerDimension; yTile++ {
			if ic.Tiles[xTile][yTile] {
				ic.drawTile(xTile, yTile)
			}

		}
	}
}

func (ic *Identicon) drawTile(xTile, yTile int) {

	xStart := (xTile * (ic.Options.ImageSize / tilesPerDimension))
	xEnd := xStart + (ic.Options.ImageSize / tilesPerDimension)

	yStart := (yTile * (ic.Options.ImageSize / tilesPerDimension))
	yEnd := yStart + (ic.Options.ImageSize / tilesPerDimension)

	// fmt.Println("x", xStart, xEnd)
	// fmt.Println("y", yStart, yEnd)

	bounds := image.Rect(xStart, yStart, xEnd, yEnd)

	// @todo possibly faster to just iterate pixels and use ic.Set()
	draw.Draw(ic, bounds, &image.Uniform{ic.Color}, image.ZP, draw.Src)

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
		ic.setTileValue(i, ic.Hash[i])
	}

	// Middle
	for i = 10; i < 15; i++ {
		ic.setTileValue(i, ic.Hash[i])
	}

	// Mirror to right
	ic.mirror()
}

func (ic *Identicon) setTileValue(pos int8, b byte) {

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

	// @todo too random? custom palette?
	ic.Color = color.NRGBA{
		R: uint8(ic.Hash[15]),
		G: uint8(ic.Hash[14]),
		B: uint8(ic.Hash[13]),
		A: uint8(255),
	}
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
	hasher.Size()
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
