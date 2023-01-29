package identicon_test

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	v0 "github.com/tsdtsdtsd/identicon"
	"github.com/tsdtsdtsd/identicon/v1"
)

var identifier = "my-test-identifier"

func TestIdentifierMustNotBeEmpty(t *testing.T) {

	icon, err := identicon.New("")

	t.Run("icon is nil", func(t *testing.T) {
		assert.Nil(t, icon)
	})

	t.Run("has error", func(t *testing.T) {
		assert.Error(t, err)
	})
}

func TestBasicIdenticon(t *testing.T) {

	icon, err := identicon.New(identifier)

	t.Run("icon is not nil", func(t *testing.T) {
		assert.NotNil(t, icon)
	})

	t.Run("has no errors", func(t *testing.T) {
		assert.NoError(t, err)
	})

	t.Run("icon has default options", func(t *testing.T) {
		defaultOptions := identicon.DefaultOptions()
		assert.Equal(t, defaultOptions, icon.Options())
	})

	t.Run("icon has correct identifier set", func(t *testing.T) {
		assert.Equal(t, identifier, icon.Identifier)
	})
}

func TestWithBGColorShouldSetOption(t *testing.T) {

	red := color.NRGBA{255, 0, 0, 255}
	defaultOptions := identicon.DefaultOptions()
	defaultOptions.BGColor = red

	icon, err := identicon.New(
		identifier,
		identicon.WithBGColor(red),
	)

	assert.NotNil(t, icon)
	assert.NoError(t, err)
	assert.Equal(t, defaultOptions, icon.Options())
}

func TestWithGridResolutionShouldSetOption(t *testing.T) {

	resolution := 8
	defaultOptions := identicon.DefaultOptions()
	defaultOptions.GridResolution = resolution

	icon, err := identicon.New(
		identifier,
		identicon.WithGridResolution(resolution),
	)

	assert.NotNil(t, icon)
	assert.NoError(t, err)
	assert.Equal(t, defaultOptions, icon.Options())
}

func TestWithGridResolutionNonPositiveValueShouldBeDiscarded(t *testing.T) {

	defaultOptions := identicon.DefaultOptions()

	t.Run("zero given", func(t *testing.T) {
		resolution := 0

		icon, err := identicon.New(
			identifier,
			identicon.WithGridResolution(resolution),
		)

		assert.NotNil(t, icon)
		assert.NoError(t, err)
		assert.Equal(t, defaultOptions, icon.Options())
	})

	t.Run("negative given", func(t *testing.T) {
		resolution := -5

		icon, err := identicon.New(
			identifier,
			identicon.WithGridResolution(resolution),
		)

		assert.NotNil(t, icon)
		assert.NoError(t, err)
		assert.Equal(t, defaultOptions, icon.Options())
	})
}

func TestHashHasTheExpectedValue(t *testing.T) {

	testSet := []struct {
		identifier   string
		expectedHash string
	}{
		{
			identifier:   "my-test",
			expectedHash: "6f14f2bed3ec4e0d3db3c7f62c2d9aef6f14f2bed3ec4e0d3db3c7f62c2d9aef",
		},
		{
			identifier:   "0",
			expectedHash: "cfcd208495d565ef66e7dff9f98764dacfcd208495d565ef66e7dff9f98764da",
		},
		{
			identifier:   "my-second-test-is-a-lot-larger-than-the-first-test-i-swear",
			expectedHash: "6eb92c6a1ff075525502a0b84470debe6eb92c6a1ff075525502a0b84470debe",
		},
	}
	for _, test := range testSet {
		icon, err := identicon.New(test.identifier)

		assert.NoError(t, err)
		assert.Equal(t, test.expectedHash, icon.HashString(), fmt.Sprintf("given: %s", test.identifier))
	}
}

func TestMatrixIsCorrect(t *testing.T) {

	// TODO: test different resolutions

	testSet := []struct {
		identifier     string
		expectedMatrix [][]bool
	}{
		{
			identifier: "my-test",
			expectedMatrix: [][]bool{
				{true, false, true, true, true},
				{false, true, false, false, true},
				{true, true, false, false, true},
				{false, true, false, false, true},
				{true, false, true, true, true},
			},
		},
		{
			identifier: "0",
			expectedMatrix: [][]bool{
				{true, false, false, false, false},
				{false, false, true, true, true},
				{true, false, false, true, false},
				{false, false, true, true, true},
				{true, false, false, false, false},
			},
		},
		{
			identifier: "my-second-test-is-a-lot-larger-than-the-first-test-i-swear",
			expectedMatrix: [][]bool{
				{true, false, false, true, true},
				{false, false, true, false, true},
				{false, false, false, false, true},
				{false, false, true, false, true},
				{true, false, false, true, true},
			},
		},
	}
	for _, test := range testSet {
		icon, err := identicon.New(test.identifier)

		assert.NoError(t, err)
		assert.Equal(t, test.expectedMatrix, icon.Matrix(), fmt.Sprintf("given: %s", test.identifier))
	}
}

func ExampleNew() {
	icon, err := identicon.New("michael@example.com")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("identicon.png")
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(file, icon)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()
	// // Output:
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		identicon.New(identifier)
	}
}

func BenchmarkNewV0(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		v0.New(identifier, nil)
	}
}
