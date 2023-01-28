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
	}
	got := identicon.DefaultOptions()

	assert.Equal(t, expected, got)
}