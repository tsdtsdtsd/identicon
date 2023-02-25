package identicon

import (
	"encoding/hex"
	"errors"
	"hash"
)

// stride is the pixel byte stride (in bytes) between vertically adjacent pixels.
const stride int = 4

// Identicon defines an identicon.
type Identicon struct {
	Identifier string
	options    *Options
	matrix     [][]bool // first dimension is columns, second is rows
	hash       []byte
	image      *Image
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
	ic.computeHash()
	ic.computeMatrix()

	return ic, nil
}

func (ic *Identicon) applyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(ic)
	}
}

func (ic *Identicon) computeHash() {
	// TODO: mandatoryByteAmount is often too large; debug with resolution=11
	tileAmount := ic.options.Resolution * ic.options.Resolution
	mandatoryByteAmount := tileAmount / 2
	resolutionIsEven := tileAmount%2 == 0

	if !resolutionIsEven {
		mandatoryByteAmount += ic.options.Resolution
	}

	sum := hashSum(ic.options.Hasher, []byte(ic.Identifier))

	for len(sum) < mandatoryByteAmount {
		addSum := hashSum(ic.options.Hasher, sum)
		sum = append(addSum, sum...)
	}

	ic.hash = sum
}

func (ic *Identicon) computeMatrix() {

	matrix := make([][]bool, ic.options.Resolution)
	even := ic.options.Resolution%2 == 0
	half := int(ic.options.Resolution / 2)

	// Columns
	for col := 0; col < ic.options.Resolution; col++ {

		// Middle col
		if col > half-1 {
			if even {
				break
			}

			if col == half {
				matrix[col] = createColumn(col, ic.hash, ic.options.Resolution, false)
			}

			continue
		}

		// First half
		matrix[col] = createColumn(col, ic.hash, ic.options.Resolution, false)

		// Mirror to second half
		colMax := len(matrix) - 1
		mirroredColNum := colMax - col
		matrix[mirroredColNum] = createColumn(mirroredColNum, ic.hash, ic.options.Resolution, true)
	}

	ic.matrix = matrix
}

// Options returns the identicons options.
func (ic *Identicon) Options() *Options {
	return ic.options
}

// Matrix returns the identicons "tile map".
func (ic *Identicon) Matrix() [][]bool {
	return ic.matrix
}

// HashString returns the hexadacimal representation of the hash as string.
func (ic *Identicon) HashString() string {
	return hex.EncodeToString(ic.hash)
}

// Image generates and returns the image.
func (ic *Identicon) Image() *Image {
	if ic.image == nil {
		ic.image = newImage(ic.options, ic.hash, ic.matrix)
	}
	return ic.image
}

func hashSum(hasher hash.Hash, in []byte) []byte {
	hasher.Reset()
	hasher.Write(in)
	return hasher.Sum(nil)
}

func tileValue(colNum int, rowNum int, hash []byte, resolution int, secondHalf bool) bool {

	realColNum := colNum
	if secondHalf {
		realColNum = (resolution - colNum - 1)
	}
	pos := (realColNum * resolution) + rowNum

	return (int(hash[pos]) & 2) > 0
}
