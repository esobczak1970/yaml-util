# yaml-util

[![Go Reference](https://pkg.go.dev/badge/github.com/esobczak1970/yaml-util.svg)](https://pkg.go.dev/github.com/esobczak1970/yaml-util)

`yaml-util` is a Go module for transforming YAML into compact, expanded, or
annotated forms.

## Requirements

- Go 1.27.1 or later.
- `gopkg.in/yaml.v3`.

## Packages

### `minify`

Compacts YAML while preserving its parsed structure.

```go
minified, err := minify.Minify(input)
```

### `maxify`

Expands YAML into a more readable block-oriented form.

```go
expanded, err := maxify.Maxify(input)
expandedTwoSpaces, err := maxify.WithIndent(input, 2)
```

### `verbose`

Adds structure/type-oriented comments using configurable options.

```go
output, err := verbose.MakeVerbose(input)

opts := verbose.DefaultOptions()
output, err = verbose.MakeVerboseWithOptions(input, opts)
```

## Installation

```bash
go get github.com/esobczak1970/yaml-util
```

## Development

Use the root Makefile:

```bash
make help
make doctor
make test
make test-race
make coverage
make lint
make check
```

Run a specific package directly when needed:

```bash
go test ./minify
go test ./maxify
go test ./verbose
```

## Repository layout

```text
yaml-util/
├── maxify/
├── minify/
├── verbose/
├── scripts/
├── Makefile
├── go.mod
└── go.sum
```

The package tests are the best executable reference for edge cases involving
anchors, aliases, scalar types, formatting, and invalid input.

## Contributing

Before submitting a change:

1. Run `make check`.
2. Add or update tests for behavior changes.
3. Keep exported APIs documented.
4. Update this README when public behavior changes.

## License

See the repository `LICENSE` file.
