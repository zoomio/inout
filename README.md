# inout

[![Build Status](https://travis-ci.com/zoomio/inout.svg?branch=master)](https://travis-ci.com/zoomio/inout)

Retrieves contents of the provided source: STDIN, HTTP(S) or FS.

## Installation

### Binary

Get the latest [release](https://github.com/zoomio/inout/releases/latest) by running this command in your shell:

For MacOS:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/inout/master/_bin/install.sh)" -o darwin
```

For MacOS (arm64):
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/inout/master/_bin/install.sh)" -o darwin arm64
```

For Linux:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/inout/master/_bin/install.sh)" -o linux
```

For Windows:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/zoomio/inout/master/_bin/install.sh)" -o windows
```

### Go dependency

```bash
go get -u github.com/zoomio/inout/...
```

## Usage

See [cmd/cli/cli.go](https://raw.githubusercontent.com/zoomio/inout/master/cmd/cli/cli.go)

## Changelog

See [CHANGELOG.md](https://raw.githubusercontent.com/zoomio/inout/master/CHANGELOG.md)

## License

Released under the [Apache License 2.0](https://raw.githubusercontent.com/zoomio/inout/master/LICENSE).
