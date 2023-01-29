package identicon

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"image"
	"image/color"
)

// Identicon defines an identicon
type Identicon struct {
	Identifier string
	options    *Options
	matrix     [][]bool // first dimension is columns, second is rows
	hash       []byte
	bounds     image.Rectangle
	// pixels holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	pixels []uint8
	// stride is the pixel stride (in bytes) between vertically adjacent pixels.
	stride int
}

var hasher = md5.New()

// var hasher = fnv.New128a()

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

	ic.initBounds()
	ic.initPixels()
	for _, opt := range opts {
		opt(ic)
	}

func (ic *Identicon) initBounds() {
	ic.bounds = image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{ic.options.ImageSize, ic.options.ImageSize},
	}
}

func (ic *Identicon) initPixels() {
	ic.pixels = make([]uint8, 4*ic.bounds.Dx()*ic.bounds.Dy())
}
	// TODO: mandatoryByteAmount is often too large; debug with resolution=11
	tileAmount := ic.options.GridResolution * ic.options.GridResolution
	mandatoryByteAmount := tileAmount / 2
	even := tileAmount%2 == 0

	// fmt.Println(tileAmount, mandatoryByteAmount)
	if !even {
		mandatoryByteAmount += ic.options.GridResolution
	}
	// fmt.Println(tileAmount, mandatoryByteAmount)
	sum := hashSum([]byte(identifier))

	for len(sum) < mandatoryByteAmount {
		sum = append(sum, sum...)
	}

	// fmt.Println(len(sum), mandatoryByteAmount)
	ic.hash = sum
	ic.matrix = computeMatrix(sum, ic.options.GridResolution)

	return ic, nil
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

func hashSum(in []byte) []byte {
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

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (ic *Identicon) pixelOffset(x, y int) int {
	return (y-ic.bounds.Min.Y)*ic.stride + (x-ic.bounds.Min.X)*4
}

func computeColumn(colNum int, hash []byte, resolution int, secondHalf bool) []bool {

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
