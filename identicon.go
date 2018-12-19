package identicon

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
)

const (
	imgWidth  = 100
	imgHeight = 100
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
func (ic *Identicon) GenerateImage() *image.Image {

	ic.createTiles()
	// ic.DebugPrintTiles()

	// Color
	colorIdx := int(ic.Hash[15])
	if colorIdx > 215 {
		colorIdx = 20
	}

	ic.Color = palette.WebSafe[colorIdx]

	// New image
	bounds := image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{imgWidth, imgHeight}}
	palette := color.Palette{ic.Color, rgb(240, 240, 240)}
	img := image.NewPaletted(bounds, palette)

	// Draw
	for xTile := 0; xTile < 5; xTile++ {
		for yTile := 0; yTile < 5; yTile++ {
			if ic.Tiles[xTile][yTile] {
				ic.drawTile(img, xTile, yTile)
			}

		}
	}

	i := img.SubImage(bounds)
	return &i
}

func (ic *Identicon) drawTile(img *image.Paletted, xTile, yTile int) *image.Paletted {

	xStart := (xTile * (imgWidth / 5))
	if xStart < 0 {
		xStart = 0
	}
	xEnd := xStart + (imgWidth / 5) - 1

	yStart := (yTile * (imgHeight / 5))
	if yStart < 0 {
		yStart = 0
	}
	yEnd := yStart + (imgHeight / 5) - 1

	// fmt.Println("x", xStart, xEnd)
	// fmt.Println("y", yStart, yEnd)

	for x := xStart; x <= xEnd; x++ {
		for y := yStart; y <= yEnd; y++ {
			img.SetColorIndex(x, y, 1)
		}
	}

	return img
}

func (ic *Identicon) createTiles() {

	tiles := make([][]bool, 5)
	for i := range tiles {
		tiles[i] = make([]bool, 5)
	}

	ic.Tiles = tiles

	// First 15 bytes of hash define pixels:
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

// DebugPrintTiles prints the tiles at positions x,y
func (ic *Identicon) DebugPrintTiles() {
	for x := range ic.Tiles {
		for y, v := range ic.Tiles[x] {
			fmt.Printf("%d,%d : %v\n", x, y, v)
		}
	}
}

func (ic *Identicon) calcTile(pos int8, b byte) {

	var lever = float32(b)
	for lever >= 2 {
		lever = lever / 3
	}

	var value bool
	if int8(lever) == 1 {
		value = true
	}

	x, y := posToXY(pos)
	ic.Tiles[x][y] = value
}

func (ic *Identicon) mirror() {
	for i := 0; i <= 1; i++ {
		x := 4 - i
		for y := 0; y < 5; y++ {
			ic.Tiles[x][y] = ic.Tiles[i][y]
			// fmt.Println("###", x, y, "|", i, y, ":", ic.Tiles[i][y])
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
