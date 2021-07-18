# mynerds

A template of the go-app project

[![PkgGoDev](https://pkg.go.dev/badge/kyoh86/mynerds)](https://pkg.go.dev/kyoh86/mynerds)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/mynerds)](https://goreportcard.com/report/github.com/kyoh86/mynerds)
[![Coverage Status](https://img.shields.io/codecov/c/github/kyoh86/mynerds.svg)](https://codecov.io/gh/kyoh86/mynerds)
[![Release](https://github.com/kyoh86/mynerds/workflows/Release/badge.svg)](https://github.com/kyoh86/mynerds/releases)

## Description

```console
$ mynerds man
```

`mynerds` provides a template of the go-app project.

## Install

### For Golang developers

```console
$ go get github.com/kyoh86/mynerds/cmd/mynerds
```

### Homebrew/Linuxbrew

```console
$ brew tap kyoh86/tap
$ brew update
$ brew install kyoh86/tap/mynerds
```

### Makepkg

```console
$ mkdir -p mynerds_build && \
  cd mynerds_build && \
  curl -iL --fail --silent https://github.com/kyoh86/mynerds/releases/latest/download/mynerds_PKGBUILD.tar.gz | tar -xvz
$ makepkg -i
```

## Available commands

Use `mynerds [command] --help` for more information about a command.
Or see the manual in [usage/mynerds.md](./usage/mynerds.md).

## Commands

Manual: [usage/mynerds.md](./usage/mynerds.md).

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This software is released under the [MIT License](http://www.opensource.org/licenses/MIT), see LICENSE.
