# Identicon

> A Go library for Identicon generation

![Latest Release Version][shields-version-img]
[![Godoc][godoc-image]][godoc-url]
![Build Status](https://github.com/tsdtsdtsd/identicon/actions/workflows/ci.yml/badge.svg?branch=rewrite)
[![Go Report Card][grc-image]][grc-url]
[![codecov][codecov-image]][codecov-url]
[![CodeQL](https://github.com/tsdtsdtsd/identicon/actions/workflows/codeql-analysis.yml/badge.svg?branch=master)](https://github.com/tsdtsdtsd/identicon/actions/workflows/codeql-analysis.yml)

This Go library helps to generate deterministic [Identicons][identicon-wiki] from strings.

---

⚠️ I'm working on a complete rewrite for v1.0 ⚠️  
It will definitely break the API and most probably also the internal hashing algorithm, so identicons generated after 1.0 will look different than before. 

---

<!-- Markdown link & img dfn's -->
[grc-image]: https://goreportcard.com/badge/github.com/tsdtsdtsd/identicon
[grc-url]: https://goreportcard.com/report/github.com/tsdtsdtsd/identicon
[godoc-image]: https://pkg.go.dev/badge/github.com/tsdtsdtsd/identicon.svg
[godoc-url]: https://pkg.go.dev/github.com/tsdtsdtsd/identicon
[codecov-image]: https://codecov.io/gh/tsdtsdtsd/identicon/branch/rewrite/graph/badge.svg
[codecov-url]: https://codecov.io/gh/tsdtsdtsd/identicon/tree/rewrite
[shields-version-img]: https://img.shields.io/github/v/release/tsdtsdtsd/identicon
[identicon-wiki]: https://en.wikipedia.org/wiki/Identicon
