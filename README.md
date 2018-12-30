# Identicon

> An experimental Go library for Identicon generation

[![Godoc][godoc-image]][godoc-url]
[![Build Status][travis-image]][travis-url]
[![Go Report Card][grc-image]][grc-url]
[![codecov][codecov-image]][codecov-url]

This Go library helps generating deterministic [Identicons][identicon-wiki], like these:

![Example](example/images/unknown.png "Example") . ![Example](example/images/test-string.png "Example") . ![Example](example/images/Amazatron3000.png "Example") . ![Example](example/images/yay-identicons.png "Example") . ![Example](example/images/m.jackson.png "Example")

![Example](example/images/12monkeys.png "Example") . ![Example](example/images/Stan.Lee.png "Example") . ![Example](example/images/gogopher.png "Example") . ![Example](example/images/notblue.png "Example") . ![Example](example/images/test.png "Example")

## Installation

```sh
go get github.com/tsdtsdtsd/identicon
```

## Usage example

Take a look at the `/example` folder, it contains a simple usage example.

Basically, you construct a new Identicon with `New`, giving it your identification string and some optional `Options`. 
The resulting struct implements `image.Image`, so you can use it flexibly.

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