package minify

// minify.go
// YAML minifier

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

// Use yaml.Node to preserve order where necessary
// When handling strings, you need to ensure proper escaping and quoting, especially if the string contains characters that could be misinterpreted in YAML (like colons or braces).
// Make sure to handle arrays properly
// write unit tests that cover things that get cleaned up this way so we can confirm what is cleaned up
// Tests should verify all minifications are valid
// In the case of tests that involve mappings the unordered map should be unmarshaled and then marshaled to ensure it is the same
// if scalar values don't get unquoted then write a function to perform this cleanup
// write a function to minify boolean from true to y, false to n and null to ~
// minify sequences in string
// minify mappings in string to inline
// minify mappings in string to remove unnecessary quotes
// minify using anchor and alias
// minify key: value to key:value spacing
// minify by removing any other whitespace or new lines that are not necessary
// could support a safe mode that would minify only as far as possible by trying each minification then validating that the original and current unmarshalled versions are valid
// code should be able to handle complex combinations of data types
