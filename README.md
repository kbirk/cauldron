# cauldron

> Eye of newt, and toe of frog

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/github.com/kbirk/cauldron)
[![Build Status](https://travis-ci.org/kbirk/cauldron.svg?branch=master)](https://travis-ci.org/kbirk/cauldron)
[![Go Report Card](https://goreportcard.com/badge/github.com/kbirk/cauldron)](https://goreportcard.com/report/github.com/kbirk/cauldron)

## Description

A little application to toil around with 2D particle effects.

## Dependencies

* [Golang](https://golang.org/):
    * Requires 1.6+ binaries are required with the `GOPATH` environment variable specified and `$GOPATH/bin` in your `PATH`.
* [go-gl/gl](https://github.com/go-gl/gl/master/README.md):
    * A cgo compiler (typically gcc).
    * On Ubuntu/Debian-based systems, the `libgl1-mesa-dev` package.
* [go-gl/glfw](https://raw.githubusercontent.com/go-gl/glfw/master/README.md):
    * On OS X, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
    * On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev` and `xorg-dev` packages.
    * On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel` packages.
    * See [here](http://www.glfw.org/docs/latest/compile.html#compile_deps) for full details.

## Installation

Clone the repository:

```bash
mkdir -p $GOPATH/src/github.com/kbirk
cd $GOPATH/src/github.com/kbirk
git clone git@github.com:kbirk/cauldron.git
```

Install dependencies:

```bash
cd cauldron
make install
```

## Usage

Build and run the executable:

```bash
go run main.go
```
