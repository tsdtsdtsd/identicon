package identicon

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

// Identicon defines an identicon
type Identicon struct {
	Identifier string
	options    *Options
	matrix     [][]bool // first dimension is columns, second is rows
	hash       []byte
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

	// Apply options
	for _, opt := range opts {
		opt(ic)
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

func computeMatrix(hash []byte, resolution int) [][]bool {

	matrix := make([][]bool, resolution)
	even := resolution%2 == 0
	half := int(resolution / 2)

	// Columns
	for col := 0; col < resolution; col++ {

		// Middle col
		if col > half-1 {
			if even {
				break
			}

			if col == half {
				matrix[col] = computeColumn(col, hash, resolution, false)
			}

			continue
		}

		// First half
		matrix[col] = computeColumn(col, hash, resolution, false)

		// Replicate to second half
		colMax := len(matrix) - 1
		mirroredColNum := colMax - col
		matrix[mirroredColNum] = computeColumn(mirroredColNum, hash, resolution, true)
	}

	return matrix
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
