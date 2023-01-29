package identicon_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tsdtsdtsd/identicon/v1"
)

func TestDefaultOptions(t *testing.T) {

	// TODO: remove this test after initial stuff, it makes no sense

	expected := &identicon.Options{
		BGColor:        color.NRGBA{240, 240, 240, 255},
		GridResolution: 5,
		ImageSize:      120,
	}
	got := identicon.DefaultOptions()

	assert.Equal(t, expected, got)
}

func TestWithBGColorOption(t *testing.T) {

	red := color.NRGBA{255, 0, 0, 255}
	optionFunc := identicon.WithBGColor(red)
	icon, err := identicon.New("id")

	optionFunc(icon)
	got := icon.Options().BGColor

	assert.NoError(t, err)
	assert.Equal(t, red, got)
}

func TestWithImageSizeOption(t *testing.T) {

	size := 200
	optionFunc := identicon.WithImageSize(size)
	icon, err := identicon.New("id")

	optionFunc(icon)
	got := icon.Options().ImageSize

	assert.NoError(t, err)
	assert.Equal(t, size, got)
}
