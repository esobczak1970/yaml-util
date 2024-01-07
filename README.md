# yaml-util

A Go module for YAML processing, providing functionalities to minify, maxify, and make YAML verbose.

## Installation

```bash
go get github.com/esobczak1970/yaml-util
```

## Usage

### Minify YAML

To minify a YAML string:

```go
package main

import (
    "fmt"
    "github.com/esobczak1970/yaml-util/minify"
)

func main() {
    yamlContent := `
    list:
      - item1
      - item2
    `

    minified, err := minify.Minify(yamlContent)
    if err != nil {
        fmt.Println("Error minifying YAML:", err)
        return
    }

    fmt.Println(minified)
}
```

This will convert the YAML content into its most compact form.

### Maxify and Verbose

Similarly, you can use the `maxify` and `verbose` packages to pretty-print or expand YAML content. (Implementations for these functions should be created following a similar pattern.)

Remember to update the README with actual examples once `maxify` and `verbose` functionalities are fully implemented.
