package identicon

import (
	"image"
	"image/color"
	"image/draw"
)

type Image struct {
	bounds  image.Rectangle
	pixels  []uint8
	options *Options
}

func newImage(options *Options, hash []byte, matrix [][]bool) *Image {
	img := Image{
		bounds: image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{options.ImageSize, options.ImageSize},
		},
		options: options,
	}

	img.pixels = make([]uint8, 4*img.bounds.Dx()*img.bounds.Dy())

	img.draw(hash, matrix)
	return &img
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (img *Image) At(x, y int) color.Color {
	return img.NRGBAAt(x, y)
}

// NRGBAAt returns the color of the pixel at (x, y) as color.NRGBA.
func (img *Image) NRGBAAt(x, y int) color.NRGBA {
	if !(image.Point{x, y}.In(img.bounds)) {
		return color.NRGBA{}
	}
	i := img.pixelOffset(x, y)
	return color.NRGBA{img.pixels[i+0], img.pixels[i+1], img.pixels[i+2], img.pixels[i+3]}
}

// Bounds returns the domain for which At can return non-zero color.
func (img *Image) Bounds() image.Rectangle {
	return img.bounds
}

// ColorModel returns the Image's color model.
func (img *Image) ColorModel() color.Model {
	return color.NRGBAModel
}

// PixOffset returns the index of the first element of pixels that corresponds to
// the pixel at (x, y).
func (img *Image) pixelOffset(x, y int) int {
	return (y-img.bounds.Min.Y)*stride*img.options.ImageSize + (x-img.bounds.Min.X)*4
}

// Set stores given color at position (x, y).
func (img *Image) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(img.bounds)) {
		return
	}

	i := img.pixelOffset(x, y)
	cNRGBA := img.ColorModel().Convert(c).(color.NRGBA)

	img.pixels[i+0] = cNRGBA.R
	img.pixels[i+1] = cNRGBA.G
	img.pixels[i+2] = cNRGBA.B
	img.pixels[i+3] = cNRGBA.A
}

func (img *Image) draw(hash []byte, matrix [][]bool) {

	fgColor := img.options.FGColor

	if fgColor == nil {
		// Last 3 bytes of hash are the RGB values
		// TODO: too random? custom palette?
		fgColor = color.NRGBA{
			R: uint8(hash[1]),
			G: uint8(hash[2]),
			B: uint8(hash[3]),
			A: uint8(255),
		}
	}

	// Background fill
	draw.Draw(img, img.Bounds(), &image.Uniform{img.options.BGColor}, image.Point{}, draw.Src)

	// Iterate tiles and draw
	for colOffset, col := range matrix {
		for rowOffset, tileIsSet := range col {
			if tileIsSet {
				img.drawTile(colOffset, rowOffset, fgColor)
			}
		}
	}
}

func (img *Image) drawTile(colOffset, rowOffset int, fgColor color.Color) {

	colStart := (colOffset * (img.options.ImageSize / img.options.Resolution))
	colEnd := colStart + (img.options.ImageSize / img.options.Resolution)

	rowStart := (rowOffset * (img.options.ImageSize / img.options.Resolution))
	rowEnd := rowStart + (img.options.ImageSize / img.options.Resolution)

	draw.Draw(
		img,
		image.Rect(colStart, rowStart, colEnd, rowEnd),
		&image.Uniform{fgColor},
		image.Point{},
		draw.Src,
	)
}

func createColumn(colNum int, hash []byte, resolution int, secondHalf bool) []bool {

	col := make([]bool, resolution)

	for rowNum := 0; rowNum < resolution; rowNum++ {
		col[rowNum] = tileValue(colNum, rowNum, hash, resolution, secondHalf)
	}

	return col
}
