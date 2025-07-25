# `enw`

A flexible Go module for sourcing configuration variables from multiple, prioritized locations.

This library simplifies configuration management by providing a single, consistent interface to find and search for
values across different environments, from local development to deployments in Kubernetes. Taste it! :heart:

<div>
  <a href="https://github.com/therenotomorrow/enw/releases" target="_blank">
    <img src="https://img.shields.io/github/v/release/therenotomorrow/enw?color=FBC02D" alt="GitHub releases">
  </a>
  <a href="https://go.dev/doc/go1.24" target="_blank">
    <img src="https://img.shields.io/badge/Go-%3E%3D%201.24-blue.svg" alt="Go 1.24">
  </a>
  <a href="https://pkg.go.dev/github.com/therenotomorrow/enw" target="_blank">
    <img src="https://godoc.org/github.com/therenotomorrow/enw?status.svg" alt="Go reference">
  </a>
  <a href="https://github.com/therenotomorrow/enw/blob/master/LICENSE" target="_blank">
    <img src="https://img.shields.io/github/license/therenotomorrow/enw?color=388E3C" alt="License">
  </a>
  <a href="https://github.com/therenotomorrow/enw/actions/workflows/ci.yml" target="_blank">
    <img src="https://github.com/therenotomorrow/enw/actions/workflows/ci.yml/badge.svg" alt="ci status">
  </a>
  <a href="https://goreportcard.com/report/github.com/therenotomorrow/enw" target="_blank">
    <img src="https://goreportcard.com/badge/github.com/therenotomorrow/enw" alt="Go report">
  </a>
  <a href="https://codecov.io/gh/therenotomorrow/enw" target="_blank">
    <img src="https://img.shields.io/codecov/c/github/therenotomorrow/enw?color=546E7A" alt="Codecov">
  </a>
</div>

## Installation

```shell
go get github.com/therenotomorrow/enw@latest
```

## Development

### System Requirements

```shell
go version
# go version go1.24.0

just --version
# just 1.40.0
```

### Download sources

```shell
PROJECT_ROOT=enw
git clone https://github.com/therenotomorrow/enw.git "$PROJECT_ROOT"
cd "$PROJECT_ROOT"
```

### Setup dependencies

```shell
# install dependencies
go mod download
go mod verify

# check code integrity
just code test smoke

# setup safe development (optional)
git config --local core.hooksPath .githooks
```

## Testing

```shell
# run quick checks
just test smoke # or just test

# run with coverage
just test cover
```

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the terms of the MIT license.
