package identicon

import (
	"image/color"
	"image/png"
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
		Debug:           true,
	})

	f, err := os.Open("example/images/" + id + ".png")
	if err != nil {
		panic("Could not open proof file: " + err.Error())
	}
	defer f.Close()

	proof, err := png.Decode(f)
	if err != nil {
		panic("Could not decode proof file: " + err.Error())
	}

	if ic.Bounds() != proof.Bounds() {
		t.Error("Generated image dimensions differ from proof file")
	}

	var diff bool

LOOP:
	for x := 0; x < proof.Bounds().Dx(); x++ {
		for y := 0; y < proof.Bounds().Dy(); y++ {

			genR, genG, genB, genA := ic.At(x, y).RGBA()
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

func TestOutOfBounds(t *testing.T) {
	ic, _ := New(id, &Options{
		BackgroundColor: RGB(235, 235, 235),
		ImageSize:       100,
	})

	blank := color.NRGBA{}
	if ic.NRGBAAt(-1, -1) != blank || ic.NRGBAAt(101, 101) != blank {
		t.Error("OOB NRGBAAt")
	}

	for _, val := range ic.Pix {
		if val == 0 {
			t.Error("OOB Set test incomplete, there already are black pixels")
			break
		}
	}

	ic.Set(-1, -1, blank)
	ic.Set(101, 101, blank)

	for _, val := range ic.Pix {
		if val == 0 {
			t.Error("OOB Set should not succeed, but did")
			break
		}
	}
}
