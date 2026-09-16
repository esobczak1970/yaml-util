# yaml-util

A comprehensive Go module for YAML processing, providing utilities to minify, maxify (expand), and add verbose comments to YAML documents.

[![Go Reference](https://pkg.go.dev/badge/github.com/esobczak1970/yaml-util.svg)](https://pkg.go.dev/github.com/esobczak1970/yaml-util)

## Features

- **Minify**: Compress YAML to its most compact valid form
- **Maxify**: Expand YAML for maximum readability
- **Verbose**: Add helpful comments and type annotations
- Preserves YAML semantics and structure
- Handles anchors, aliases, and complex nested structures
- Comprehensive error handling
- Well-tested with extensive unit tests

## Installation

```bash
go get github.com/esobczak1970/yaml-util
```

## Usage

### Minify YAML

Compress YAML to its most compact form while preserving validity:

```go
package main

import (
    "fmt"
    "github.com/esobczak1970/yaml-util/minify"
)

func main() {
    yamlContent := `
config:
  database:
    host: localhost
    port: 5432
  features:
    - auth
    - api
    - frontend
`

    minified, err := minify.Minify(yamlContent)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println(minified)
    // Output:
    // config: {database: {host: localhost, port: 5432}, features: [auth, api, frontend]}
}
```

**Minify Features:**
- Inline short mappings: `{key:value, key2:value2}`
- Inline short sequences: `[item1, item2, item3]`
- Boolean shorthand: `true → y`, `false → n`, `null → ~`
- Remove unnecessary whitespace
- Preserve anchors & aliases: `&ref`, `*ref`

### Maxify YAML

Expand YAML for better readability:

```go
package main

import (
    "fmt"
    "github.com/esobczak1970/yaml-util/maxify"
)

func main() {
    compactYAML := "config:{host:localhost,port:5432,debug:y}"

    expanded, err := maxify.Maxify(compactYAML)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println(expanded)
    // Output:
    // config:
    //     host: localhost
    //     port: 5432
    //     debug: true
}
```

**Custom Indentation:**

```go
// Use 2-space indentation
expanded, err := maxify.WithIndent(yamlContent, 2)
```

**Maxify Features:**
- Expand inline mappings to block style
- Expand inline sequences to block style
- Convert shorthand booleans: `y → true`, `n → false`, `~ → null`
- Consistent indentation (default 4 spaces, configurable)
- Proper spacing: `key: value`

### Verbose YAML

Add helpful comments and type annotations:

```go
package main

import (
    "fmt"
    "github.com/esobczak1970/yaml-util/verbose"
)

func main() {
    yamlContent := `
app:
  name: myapp
  port: 8080
  debug: true
  features:
    - auth
    - api
`

    verboseYAML, err := verbose.MakeVerbose(yamlContent)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println(verboseYAML)
    // Output:
    // # YAML Document
    // # Mapping with 1 key-value pairs
    // app: # Nested mapping (4 keys)
    //     name: myapp # string
    //     port: 8080 # integer
    //     debug: true # boolean
    //     features: # List with 2 items
    //         - auth # string
    //         - api # string
}
```

**Custom Options:**

```go
opts := verbose.Options{
    AddTypeComments:      true,  // Add type annotations (string, integer, etc.)
    AddStructureComments: true,  // Add structural information (mappings, lists)
    AddExamples:          false, // Add example values (future feature)
    Indent:               4,     // Indentation level
}

verboseYAML, err := verbose.MakeVerboseWithOptions(yamlContent, opts)
```

**Verbose Features:**
- Type detection: `# string`, `# integer`, `# boolean`, `# float`, `# null`
- Structure comments: `# Mapping with N keys`, `# List with N items`
- Anchor annotations: `# anchor: defaults`
- Alias references: `# Alias reference`
- Preserves all YAML semantics

## API Reference

### minify Package

```go
// Minify reduces YAML to its most compact valid form
func Minify(inputYAML string) (string, error)
```

### maxify Package

```go
// Maxify expands YAML into a human-readable format
func Maxify(inputYAML string) (string, error)

// WithIndent expands YAML with custom indentation
func WithIndent(inputYAML string, indent int) (string, error)
```

### verbose Package

```go
// MakeVerbose adds comments and structure hints to YAML
func MakeVerbose(inputYAML string) (string, error)

// MakeVerboseWithOptions adds comments with custom options
func MakeVerboseWithOptions(inputYAML string, opts Options) (string, error)

// DefaultOptions returns the default verbose options
func DefaultOptions() Options

// Options controls comment generation behavior
type Options struct {
    AddTypeComments      bool // Add type annotations
    AddStructureComments bool // Add structural information
    AddExamples          bool // Add example values
    Indent               int  // Indentation level
}
```

## Examples

### Round-trip Processing

```go
// Minify -> Maxify -> Verbose workflow
original := `
config:
  database:
    host: localhost
    port: 5432
`

// Minify for storage/transmission
minified, _ := minify.Minify(original)
// Result: config:{database:{host:localhost,port:5432}}

// Maxify for readability
expanded, _ := maxify.Maxify(minified)
// Result: properly indented block format

// Verbose for documentation
documented, _ := verbose.MakeVerbose(expanded)
// Result: includes type and structure comments
```

### Configuration File Processing

```go
// Read config file
data, _ := os.ReadFile("config.yaml")

// Minify for deployment
minified, _ := minify.Minify(string(data))
os.WriteFile("config.min.yaml", []byte(minified), 0644)

// Create verbose version for documentation
verbose, _ := verbose.MakeVerbose(string(data))
os.WriteFile("config.documented.yaml", []byte(verbose), 0644)
```

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run specific package tests:

```bash
go test ./minify
go test ./maxify
go test ./verbose
```

Run with verbose output:

```bash
go test -v ./...
```

## Project Structure

```
yaml-util/
├── minify/
│   ├── minify.go       # Minification implementation
│   └── minify_test.go  # Comprehensive tests (100+ test cases)
├── maxify/
│   ├── maxify.go       # Maxification implementation
│   └── maxify_test.go  # Comprehensive tests (80+ test cases)
├── verbose/
│   ├── verbose.go      # Verbose implementation
│   └── verbose_test.go # Comprehensive tests (90+ test cases)
├── scripts/
│   ├── Make.sh         # Setup script
│   └── Test.sh         # Testing script
├── go.mod
├── go.sum
└── README.md
```

## Requirements

- Go 1.25 or higher
- gopkg.in/yaml.v3

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- New features include tests
- Code follows Go conventions
- Documentation is updated

## License

[Specify your license here]

## Acknowledgments

Built with [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)
