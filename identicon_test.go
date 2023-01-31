package identicon_test

import (
	"fmt"
	"hash/fnv"
	"image/color"
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tsdtsdtsd/identicon"
)

var (
	identifier = "my-test-identifier"
	hasher     = fnv.New128()
)

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

	icon, err := identicon.New(identifier, identicon.WithHasher(hasher))

	t.Run("icon is not nil", func(t *testing.T) {
		assert.NotNil(t, icon)
	})

	t.Run("has no errors", func(t *testing.T) {
		assert.NoError(t, err)
	})

	t.Run("icon has default options", func(t *testing.T) {
		defaultOptions := identicon.DefaultOptions()
		defaultOptions.Hasher = hasher
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
	defaultOptions.Hasher = hasher

	icon, err := identicon.New(
		identifier,
		identicon.WithBGColor(red),
		identicon.WithHasher(hasher),
	)

	assert.NotNil(t, icon)
	assert.NoError(t, err)
	assert.Equal(t, defaultOptions, icon.Options())
}

func TestWithGridResolutionShouldSetOption(t *testing.T) {

	resolution := 8
	defaultOptions := identicon.DefaultOptions()
	defaultOptions.GridResolution = resolution
	defaultOptions.Hasher = hasher

	icon, err := identicon.New(
		identifier,
		identicon.WithGridResolution(resolution),
		identicon.WithHasher(hasher),
	)

	assert.NotNil(t, icon)
	assert.NoError(t, err)
	assert.Equal(t, defaultOptions, icon.Options())
}

func TestWithGridResolutionNonPositiveValueShouldBeDiscarded(t *testing.T) {

	defaultOptions := identicon.DefaultOptions()
	defaultOptions.Hasher = hasher

	t.Run("zero given", func(t *testing.T) {
		resolution := 0

		icon, err := identicon.New(
			identifier,
			identicon.WithGridResolution(resolution),
			identicon.WithHasher(hasher),
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
			identicon.WithHasher(hasher),
		)

		assert.NotNil(t, icon)
		assert.NoError(t, err)
		assert.Equal(t, defaultOptions, icon.Options())
	})
}

func TestWithImageSizeShouldSetOption(t *testing.T) {

	size := 60
	defaultOptions := identicon.DefaultOptions()
	defaultOptions.ImageSize = size
	defaultOptions.Hasher = hasher

	icon, err := identicon.New(
		identifier,
		identicon.WithImageSize(size),
		identicon.WithHasher(hasher),
	)

	assert.NotNil(t, icon)
	assert.NoError(t, err)
	assert.Equal(t, defaultOptions, icon.Options())
}

func TestWithImageSizeNonPositiveValueShouldBeDiscarded(t *testing.T) {

	defaultOptions := identicon.DefaultOptions()
	defaultOptions.Hasher = hasher

	t.Run("zero given", func(t *testing.T) {
		size := 0

		icon, err := identicon.New(
			identifier,
			identicon.WithImageSize(size),
			identicon.WithHasher(hasher),
		)

		assert.NotNil(t, icon)
		assert.NoError(t, err)
		assert.Equal(t, defaultOptions, icon.Options())
	})

	t.Run("negative given", func(t *testing.T) {
		size := -5

		icon, err := identicon.New(
			identifier,
			identicon.WithImageSize(size),
			identicon.WithHasher(hasher),
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
			expectedHash: "e25338e76f1ce3c59da92abae869587cd065bb972f4ff78c11f12d9a75f9806c",
		},
		{
			identifier:   "0",
			expectedHash: "35fbae7d37875a800d8614f09817791dd228cb69101a8caf78912b704e4a144f",
		},
		{
			identifier:   "my-second-test-is-a-lot-larger-than-the-first-test-i-swear",
			expectedHash: "6fa315b3bd5a894a8215f10eca7146c90774be2442906c4540e22e04ee99b649",
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
				{true, true, false, true, true},
				{false, true, false, false, false},
				{true, true, false, false, false},
				{false, true, false, false, false},
				{true, true, false, true, true},
			},
		},
		{
			identifier: "0",
			expectedMatrix: [][]bool{
				{false, true, true, false, true},
				{true, true, false, false, true},
				{false, false, false, true, false},
				{true, true, false, false, true},
				{false, true, true, false, true},
			},
		},
		{
			identifier: "my-second-test-is-a-lot-larger-than-the-first-test-i-swear",
			expectedMatrix: [][]bool{
				{true, true, false, true, false},
				{true, false, true, true, false},
				{false, true, true, false, true},
				{true, false, true, true, false},
				{true, true, false, true, false},
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
	icon, err := identicon.New("michael@example.com", identicon.WithBGColor(color.RGBA{220, 220, 220, 255}))
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
	// Output:
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		identicon.New(identifier)
	}
}
