# utils-go

[![godoc](https://pkg.go.dev/badge/oss.terrastruct.com/utils-go.svg)](https://pkg.go.dev/oss.terrastruct.com/utils-go)
[![ci](https://github.com/terrastruct/utils-go/actions/workflows/ci.yml/badge.svg)](https://github.com/terrastruct/utils-go/actions/workflows/ci.yml)
[![license](https://img.shields.io/github/license/terrastruct/utils-go?color=9cf)](./LICENSE.txt)

Terrastruct's general purpose go libraries.

See https://pkg.go.dev/oss.terrastruct.com/utils-go for docs.

If there's enough external demand for a single package to be split off into its
own repo from this collection we will. Feel free to open an issue to request.

<!-- toc -->

* [assert](#assert)
* [diff](#diff)
* [xdefer](#xdefer)
* [cmdlog](#cmdlog)
* [xos](#xos)
* [xrand](#xrand)
* [xcontext](#xcontext)
* [xjson](#xjson)
* [xfmt](#xfmt)

<!-- tocstop -->

## Package Index

godoc is the canonical reference but we've provided this index as the godoc UI is frankly
garbage after the move to pkg.go.dev. It's nowhere near as clear and responsive as the old
UI. If this feedback reaches the authors of pkg.go.dev, please revert the UI back to what
it was with godoc.org.

### [assert](./assert)

assert provides test assertion helpers.

### [diff](./diff)

diff providers functions to diff strings, files and general Go values with git diff.

### [xdefer](./xdefer)

xdefer annotates all errors returned from a function transparently.

### [cmdlog](./cmdlog)

cmdlog implements color leveled logging for command line tools.

### [xos](./xos)

xos provides OS helpers.

### [xrand](./xrand)

xrand provides helpers for generating useful random values.
We use it mainly for generating inputs to tests.

### [xcontext](./xcontext)

xcontext implements indispensable context helpers.

### [xjson](xjson)

xjson implements basic JSON helpers.

### [xfmt](xfmt)

xfmt provides formatting helpers.
