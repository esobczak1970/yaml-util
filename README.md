# yaml-util

A Go module for YAML processing, providing functionalities to minify, maxify, and make YAML verbose.

## Installation

```bash
go get github.com/esobczak1970/yaml-util
Usage
Minify YAML
To minify a YAML string:

go
Copy
Edit
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
This will convert the YAML content into its most compact form, removing unnecessary whitespace while preserving structure.

Features
✅ Inline Mappings - Converts nested mappings into {key:value, key2:value2} format.
✅ Boolean & Null Shortening - true → y, false → n, null → ~.
✅ Preserves Anchors & Aliases - &default and *default remain intact.
✅ Handles Nested Lists & Maps - Compact formatting of structures.

Maxify YAML (TODO)
📌 Planned Feature: maxify will expand YAML for better readability by formatting everything in a human-friendly structure.

Example Usage (Coming Soon)
go
Copy
Edit
package main

import (
    "fmt"
    "github.com/esobczak1970/yaml-util/maxify"
)

func main() {
    yamlContent := "key:value"

    maxified, err := maxify.Maxify(yamlContent)
    if err != nil {
        fmt.Println("Error maxifying YAML:", err)
        return
    }

    fmt.Println(maxified)
}
📌 Expected Output:

makefile
Copy
Edit
key: value
Verbose YAML (TODO)
📌 Planned Feature: verbose will add comments and structure hints to YAML.

Example Usage (Coming Soon)
go
Copy
Edit
package main

import (
    "fmt"
    "github.com/esobczak1970/yaml-util/verbose"
)

func main() {
    yamlContent := "key:value"

    verboseYAML, err := verbose.MakeVerbose(yamlContent)
    if err != nil {
        fmt.Println("Error making YAML verbose:", err)
        return
    }

    fmt.Println(verboseYAML)
}
📌 Expected Output:

vbnet
Copy
Edit
# Key-value mapping
key: value  # Main configuration key
Development Notes
The following files are still TODO:

maxify/maxify.go
maxify/maxify_test.go
verbose/verbose.go
verbose/verbose_test.go
📌 A Make.sh script will ensure that placeholder files exist.