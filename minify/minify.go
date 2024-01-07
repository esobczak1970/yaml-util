package minify

// minify.go
// YAML minifier

// TODO:
// 1. Write a function to minify boolean values from true to 'y', false to 'n', and null to '~'.
// 2. Minify sequences in string representations.
// 3. Minify mappings in string representations to inline format.
// 4. Minify mappings in string representations to remove unnecessary quotes.
// 5. Minify using anchor and alias to reduce redundancy.
// 6. Minify key: value spacing to key:value.
// 7. Remove any other unnecessary whitespace or new lines.
// 8. Implement a 'safe mode' that minifies only as far as possible while ensuring that the unmarshaled versions of the original and current YAML are valid.
// 9. Ensure code can handle complex combinations of data types.

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

// Minify takes a YAML string, minifies it and returns the minified YAML.
func Minify(inputYAML string) (string, error) {
	var rootNode yaml.Node

	// Unmarshal the YAML into the rootNode
	err := yaml.Unmarshal([]byte(inputYAML), &rootNode)
	if err != nil {
		return "", err
	}

	// Apply minification to the rootNode
	minifyNode(&rootNode)

	// Create a buffer to hold the minified YAML
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2) // Reduced indentation for compactness

	// Encode the rootNode back to YAML and put it in the buffer
	if err := encoder.Encode(&rootNode); err != nil {
		return "", err
	}

	// Convert buffer contents to a string and trim leading whitespace
	minifiedYAML := strings.TrimLeft(buffer.String(), " \t")

	// Add a newline at the end if it's not already there
	if !strings.HasSuffix(minifiedYAML, "\n") {
		minifiedYAML += "\n"
	}

	return minifiedYAML, nil
}

func minifyNode(node *yaml.Node) {
	switch node.Kind {
	case yaml.MappingNode, yaml.SequenceNode:
		for _, n := range node.Content {
			minifyNode(n)
		}
	case yaml.ScalarNode:
		// Trim spaces and remove unnecessary quotes from scalar values
		node.Value = strings.TrimSpace(node.Value)
		if node.Style == yaml.DoubleQuotedStyle || node.Style == yaml.SingleQuotedStyle {
			node.Style = yaml.TaggedStyle
		}
	}
}
