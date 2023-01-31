package identicon

import (
	"hash"
	"hash/fnv"
	"image/color"
)

// Options define customizable settings
type Options struct {
	GridResolution int
	ImageSize      int
	Hasher         hash.Hash
	BGColor        color.Color
	FGColor        color.Color
}

// Option changes a single option
type Option func(*Identicon)

// DefaultOptions are the baseline for a new identicon
func DefaultOptions() *Options {
	return &Options{
		BGColor:        color.NRGBA{240, 240, 240, 255},
		GridResolution: 5,
		ImageSize:      100,
		Hasher:         fnv.New128(),
	}
}

// WithBGColor returns an option that sets the identicon's background color to given color.
func WithBGColor(c color.Color) Option {
	return func(i *Identicon) {
		i.options.BGColor = c
	}
}

// WithFGColor returns an option that sets the identicon's background color to given color.
func WithFGColor(c color.Color) Option {
	return func(i *Identicon) {
		i.options.FGColor = c
	}
}

// WithGridResolution returns an option that sets the identicon's grid resolution to given amount.
// The option will be discarded silently if given value is non-positive.
func WithGridResolution(resolution int) Option {
	return func(i *Identicon) {
		if resolution <= 0 {
			return
		}

		i.options.GridResolution = resolution
	}
}

// WithImageSize returns an option that sets the identicon's image size to given amount.
// The option will be discarded silently if given value is non-positive.
func WithImageSize(size int) Option {
	return func(i *Identicon) {
		if size <= 0 {
			return
		}

		i.options.ImageSize = size
	}
}

// WithHasher returns an option that sets the identicon's hash generator.
// The option will be discarded silently if nil given.
func WithHasher(hasher hash.Hash) Option {
	return func(i *Identicon) {
		if hasher == nil {
			return
		}

		i.options.Hasher = hasher
	}
}
