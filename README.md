# my-nerds

A template of the go-app project

[![PkgGoDev](https://pkg.go.dev/badge/kyoh86/my-nerds)](https://pkg.go.dev/kyoh86/my-nerds)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/my-nerds)](https://goreportcard.com/report/github.com/kyoh86/my-nerds)
[![Coverage Status](https://img.shields.io/codecov/c/github/kyoh86/my-nerds.svg)](https://codecov.io/gh/kyoh86/my-nerds)
[![Release](https://github.com/kyoh86/my-nerds/workflows/Release/badge.svg)](https://github.com/kyoh86/my-nerds/releases)

## Description

```console
$ my-nerds man
```

`my-nerds` provides a template of the go-app project.

## Install

### For Golang developers

```console
$ go get github.com/kyoh86/my-nerds/cmd/my-nerds
```

### Homebrew/Linuxbrew

```console
$ brew tap kyoh86/tap
$ brew update
$ brew install kyoh86/tap/my-nerds
```

### Makepkg

```console
$ mkdir -p my-nerds_build && \
  cd my-nerds_build && \
  curl -iL --fail --silent https://github.com/kyoh86/my-nerds/releases/latest/download/my-nerds_PKGBUILD.tar.gz | tar -xvz
$ makepkg -i
```

## Available commands

Use `my-nerds [command] --help` for more information about a command.
Or see the manual in [usage/my-nerds.md](./usage/my-nerds.md).

## Commands

Manual: [usage/my-nerds.md](./usage/my-nerds.md).

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This software is released under the [MIT License](http://www.opensource.org/licenses/MIT), see LICENSE.
