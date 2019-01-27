# Identicon

> An experimental Go library for Identicon generation

[![Godoc][godoc-image]][godoc-url]
[![Build Status][travis-image]][travis-url]
[![Go Report Card][grc-image]][grc-url]
[![codecov][codecov-image]][codecov-url]

This Go library helps generating deterministic [Identicons][identicon-wiki], like these:

![Example Banner](identicon-banner.png "Example Banner")

## Installation

```sh
go get github.com/tsdtsdtsd/identicon
```

## Usage example

Basically, you construct a new `Identicon` type with `identicon.New()`, giving it your identification string and some optional `identicon.Options`. 
The resulting struct implements `image.Image` and `draw.Image`, so you can use it flexibly.

Just import the library and create a new identicon:

```go
package main

import (
    "log"
    "github.com/tsdtsdtsd/identicon"
) 

func main() {
    
    ic, err := identicon.New(
        
        // The identicon ID string is mandatory.
	// Same string will always result in the same generated identicon.
	// Typically this is a username or email address.
        "identicon",
        
	// You can define custom options or pass nil for defaults
	&identicon.Options{
	    BackgroundColor: identicon.RGB(240, 240, 240),
	},
    )
    
    if err != nil {
	log.Fatal(err)
    }
    
    // Now you are free to use identicon `ic` as any other image.Image or draw.Image interface
    fi, _ := os.Create("my-file.png")
    png.Encode(fi, ic)
}
```

### Banner example

You can find another example in the `/example` folder. It contains an application, which generates the above image.
It also helps me to test the algorythm for changes.

<!-- Markdown link & img dfn's -->
[grc-image]: https://goreportcard.com/badge/github.com/tsdtsdtsd/identicon
[grc-url]: https://goreportcard.com/report/github.com/tsdtsdtsd/identicon
[godoc-image]: https://godoc.org/github.com/tsdtsdtsd/identicon?status.svg
[godoc-url]: https://godoc.org/github.com/tsdtsdtsd/identicon
[travis-image]: https://travis-ci.org/tsdtsdtsd/identicon.svg?branch=master
[travis-url]: https://travis-ci.org/tsdtsdtsd/identicon
[codecov-image]: https://codecov.io/gh/tsdtsdtsd/identicon/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/tsdtsdtsd/identicon
[identicon-wiki]: https://en.wikipedia.org/wiki/Identicon