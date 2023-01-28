package identicon_test

import (
	"fmt"
	"image/color"
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
	t.Skip()
	testSet := []struct {
		identifier   string
		expectedHash string
	}{
		{
			identifier:   "my-test",
			expectedHash: "4cb602ad084ff78d76a7f90aa5901b22",
		},
		{
			identifier:   "0",
			expectedHash: "d228cb69401a8caf78912b704e4a4f8f",
		},
		{
			identifier:   "my-second-test-is-a-lot-larger-than-the-first-test-i-swear",
			expectedHash: "33d9c9980e14bd067e8e5995c54270a1",
		},
	}
	for _, test := range testSet {
		icon, err := identicon.New(test.identifier)

		assert.NoError(t, err)
		assert.Equal(t, test.expectedHash, icon.HashString(), fmt.Sprintf("given: %s", test.identifier))
	}
}

func TestMatrixIsCorrect(t *testing.T) {

	t.Skip()
	// TODO: test different resolutions

	testSet := []struct {
		identifier     string
		expectedMatrix [][]bool
	}{
		{
			identifier: "my-test",
			expectedMatrix: [][]bool{
				{false, true, true, false, false},
				{true, true, false, true, true},
				{false, true, false, false, true},
				{true, true, false, true, true},
				{false, true, true, false, false},
			},
		},
		{
			identifier: "0",
			expectedMatrix: [][]bool{
				{true, false, true, false, false},
				{true, false, true, false, false},
				{true, false, true, true, true},
				{true, false, true, false, false},
				{true, false, true, false, false},
			},
		},
		{
			identifier: "my-second-test-is-a-lot-larger-than-the-first-test-i-swear",
			expectedMatrix: [][]bool{
				{true, false, false, false, true},
				{false, false, true, true, true},
				{false, false, false, true, false},
				{false, false, true, true, true},
				{true, false, false, false, true},
			},
		},
	}
	for _, test := range testSet {
		icon, err := identicon.New(test.identifier)

		assert.NoError(t, err)
		assert.Equal(t, test.expectedMatrix, icon.Matrix(), fmt.Sprintf("given: %s", test.identifier))
	}
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
