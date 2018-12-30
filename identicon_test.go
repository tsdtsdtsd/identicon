package identicon

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"testing"
)

var (
	id = "test-string"
)

func TestNew(t *testing.T) {

	ic, err := New(id, nil)
	if ic == nil {
		t.Error("New Identicon struct is nil")
	}
	if err != nil {
		t.Error("Error creating new Identicon struct ", err)
	}

	if ic.ID != id {
		t.Error("ID error: expected", id, ", got ", ic.ID)
	}

	if ic.HashString() != "661f8009fa8e56a9d0e94a0a644397d7" {
		t.Error("MD5 error: expected 661f8009fa8e56a9d0e94a0a644397d7, got ", ic.HashString())
	}
}

func TestGernerate(t *testing.T) {
	ic, _ := New(id, &Options{
		BackgroundColor: RGB(235, 235, 235),
		ImageSize:       100,
	})
	generated := ic.GenerateImage()

	f, err := os.Open(id + ".png")
	if err != nil {
		panic("Could not open proof file: " + err.Error())
	}
	defer f.Close()

	proof, err := png.Decode(f)
	if err != nil {
		panic("Could not decode proof file: " + err.Error())
	}

	if generated.Bounds() != proof.Bounds() {
		t.Error("Generated image dimensions differ from proof file")
	}

	var diff bool

LOOP:
	for x := 0; x < proof.Bounds().Dx(); x++ {
		for y := 0; y < proof.Bounds().Dy(); y++ {

			genR, genG, genB, genA := generated.At(x, y).RGBA()
			proR, proG, proB, proA := proof.At(x, y).RGBA()

			// genR, genG, genB, genA := uint8(genR32), uint8(genG32), uint8(genB32), uint8(genA32)
			// proR, proG, proB, proA := uint8(proR32), uint8(proG32), uint8(proB32), uint8(proA32)

			if genR != proR || genG != proG || genB != proB || genA != proA {
				// e := fmt.Sprintf("Compare error at %d:%d", x, y)
				// t.Error(e)
				// fmt.Println("X:Y ", x, y, " || ", genR, genG, genB, genA, "||", proR, proG, proB, proA)
				diff = true
				break LOOP
			}
		}
	}

	if diff {
		t.Error("Generated image not identical to proof file")
	}
}

// func TestColorPalette(t *testing.T) {
// 	ic, _ := New(id, nil)
// 	ic.Hash = []byte{
// 		byte(0),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(1),
// 		byte(255), // Color byte may not overflow palette size of 215
// 	}

// 	generated := ic.GenerateImage()

// 	if generated == nil {
// 		t.Error("Generation for color palette test failed")
// 	}

// }

func TestDebug(t *testing.T) {
	ic, _ := New(id, &Options{Debug: true})
	generated := ic.GenerateImage()

	if generated == nil {
		t.Error("Generation for debug test failed")
	}
}

func FastCompare(img1, img2 *image.RGBA) (int64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := int64(0)

	for i := 0; i < len(img1.Pix); i++ {
		accumError += int64(sqDiffUInt8(img1.Pix[i], img2.Pix[i]))
	}

	return int64(math.Sqrt(float64(accumError))), nil
}

func sqDiffUInt8(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}
