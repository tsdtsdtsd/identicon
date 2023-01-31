package identicon

import (
	"encoding/hex"
	"errors"
	"hash"
	"image"
	"image/color"
	"image/draw"
)

// stride is the pixel stride (in bytes) between vertically adjacent pixels.
const stride int = 4

// Identicon defines an identicon
type Identicon struct {
	Identifier string
	options    *Options
	matrix     [][]bool // first dimension is columns, second is rows
	hash       []byte
	bounds     image.Rectangle
	// pixels holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at pixels[(y-Rect.Min.Y)*stride*imgWidth + (x-Rect.Min.X)*stride].
	pixels []uint8
}

// New returns the identicon for given identifier.
// Additional options can specify the identicon.
// DefaultOptions() will be used as a baseline.
func New(identifier string, opts ...Option) (*Identicon, error) {

	if identifier == "" {
		return nil, errors.New("identifier can not be empty")
	}

	ic := &Identicon{
		Identifier: identifier,
		options:    DefaultOptions(),
	}

	ic.applyOptions(opts...)
	ic.initImage()
	ic.computeHash()
	ic.computeMatrix()
	ic.draw()

	return ic, nil
}

func (ic *Identicon) applyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(ic)
	}
}

func (ic *Identicon) initImage() {
	ic.bounds = image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{ic.options.ImageSize, ic.options.ImageSize},
	}

	ic.pixels = make([]uint8, 4*ic.bounds.Dx()*ic.bounds.Dy())
}

func (ic *Identicon) computeHash() {
	// TODO: mandatoryByteAmount is often too large; debug with resolution=11
	tileAmount := ic.options.GridResolution * ic.options.GridResolution
	mandatoryByteAmount := tileAmount / 2
	resolutionIsEven := tileAmount%2 == 0

	// fmt.Println(tileAmount, mandatoryByteAmount)
	if !resolutionIsEven {
		mandatoryByteAmount += ic.options.GridResolution
	}
	// fmt.Println(tileAmount, mandatoryByteAmount)
	sum := hashSum(ic.options.Hasher, []byte(ic.Identifier))

	// TODO: need better idea for workaround - this ends up in strange and repeated patterns if gridres is >8
	for len(sum) < mandatoryByteAmount {
		addSum := hashSum(ic.options.Hasher, sum)
		sum = append(addSum, sum...)
	}

	// fmt.Println(len(sum), mandatoryByteAmount)
	ic.hash = sum
}

func (ic *Identicon) computeMatrix() {

	matrix := make([][]bool, ic.options.GridResolution)
	even := ic.options.GridResolution%2 == 0
	half := int(ic.options.GridResolution / 2)

	// Columns
	for col := 0; col < ic.options.GridResolution; col++ {

		// Middle col
		if col > half-1 {
			if even {
				break
			}

			if col == half {
				matrix[col] = createColumn(col, ic.hash, ic.options.GridResolution, false)
			}

			continue
		}

		// First half
		matrix[col] = createColumn(col, ic.hash, ic.options.GridResolution, false)

		// Replicate to second half
		colMax := len(matrix) - 1
		mirroredColNum := colMax - col
		matrix[mirroredColNum] = createColumn(mirroredColNum, ic.hash, ic.options.GridResolution, true)
	}

	ic.matrix = matrix
}

func (ic *Identicon) draw() {

	// Last 3 bytes of hash are the RGB values
	// TODO: too random? custom palette?
	hashLength := len(ic.hash)
	fgColor := color.NRGBA{
		R: uint8(ic.hash[hashLength-1]),
		G: uint8(ic.hash[hashLength-2]),
		B: uint8(ic.hash[hashLength-3]),
		A: uint8(255),
	}

	// Background fill
	draw.Draw(ic, ic.Bounds(), &image.Uniform{ic.options.BGColor}, image.Point{}, draw.Src)

	// Iterate tiles and draw
	for colOffset, col := range ic.matrix {
		for rowOffset, tileIsSet := range col {
			if tileIsSet {
				ic.drawTile(colOffset, rowOffset, fgColor)
			}
		}
	}
}

func (ic *Identicon) drawTile(colOffset, rowOffset int, fgColor color.Color) {

	colStart := (colOffset * (ic.options.ImageSize / ic.options.GridResolution))
	colEnd := colStart + (ic.options.ImageSize / ic.options.GridResolution)

	rowStart := (rowOffset * (ic.options.ImageSize / ic.options.GridResolution))
	rowEnd := rowStart + (ic.options.ImageSize / ic.options.GridResolution)

	draw.Draw(
		ic,
		image.Rect(colStart, rowStart, colEnd, rowEnd),
		&image.Uniform{fgColor},
		image.Point{},
		draw.Src,
	)

}

// Options returns the identicons options.
func (ic *Identicon) Options() *Options {
	return ic.options
}

// Matrix
func (ic *Identicon) Matrix() [][]bool {
	return ic.matrix
}

// HashString returns the hexadacimal representation of the hash as string
func (ic *Identicon) HashString() string {
	return hex.EncodeToString(ic.hash)
}

func hashSum(hasher hash.Hash, in []byte) []byte {
	hasher.Reset()
	hasher.Write(in)
	return hasher.Sum(nil)
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (ic *Identicon) At(x, y int) color.Color {
	return ic.NRGBAAt(x, y)
}

// NRGBAAt returns the color of the pixel at (x, y) as color.NRGBA.
func (ic *Identicon) NRGBAAt(x, y int) color.NRGBA {
	if !(image.Point{x, y}.In(ic.bounds)) {
		return color.NRGBA{}
	}
	i := ic.pixelOffset(x, y)
	return color.NRGBA{ic.pixels[i+0], ic.pixels[i+1], ic.pixels[i+2], ic.pixels[i+3]}
}

// Bounds returns the domain for which At can return non-zero color.
func (ic *Identicon) Bounds() image.Rectangle {
	return ic.bounds
}

// ColorModel returns the Image's color model.
func (ic *Identicon) ColorModel() color.Model {
	return color.NRGBAModel
}

// PixOffset returns the index of the first element of pixels that corresponds to
// the pixel at (x, y).
func (ic *Identicon) pixelOffset(x, y int) int {
	return (y-ic.bounds.Min.Y)*stride*ic.options.ImageSize + (x-ic.bounds.Min.X)*4
}

// Set stores given color at position (x, y).
func (ic *Identicon) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(ic.bounds)) {
		return
	}

	i := ic.pixelOffset(x, y)
	cNRGBA := ic.ColorModel().Convert(c).(color.NRGBA)

	ic.pixels[i+0] = cNRGBA.R
	ic.pixels[i+1] = cNRGBA.G
	ic.pixels[i+2] = cNRGBA.B
	ic.pixels[i+3] = cNRGBA.A
}

func createColumn(colNum int, hash []byte, resolution int, secondHalf bool) []bool {

	col := make([]bool, resolution)

	for rowNum := 0; rowNum < resolution; rowNum++ {
		col[rowNum] = tileValue(colNum, rowNum, hash, resolution, secondHalf)
	}

	return col
}

func tileValue(colNum int, rowNum int, hash []byte, resolution int, secondHalf bool) bool {

	realColNum := colNum
	if secondHalf {
		realColNum = (resolution - colNum - 1)
	}
	pos := (realColNum * resolution) + rowNum
	// TODO: remove
	// fmt.Println(colNum, rowNum, pos, secondHalf)
	return (int(hash[pos]) & 2) > 0
}
