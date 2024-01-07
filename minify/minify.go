package minify

// minify.go
// YAML minifier

// TODO:
// Minify mappings in string representations to inline format.
// Minify mappings in string representations to remove unnecessary quotes.
// Minify using anchor and alias to reduce redundancy.
// Minify key: value spacing to key:value.
// Remove any other unnecessary whitespace or new lines.
// Implement a 'safe mode' that minifies only as far as possible while ensuring that the unmarshaled versions of the original and current YAML are valid.
// Ensure code can handle complex combinations of data types.

import (
	"bytes"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Minify takes a YAML string, minifies it and returns the minified YAML.
func Minify(inputYAML string) (string, error) {
	// Preprocess the input string
	inputYAML = preprocessMinifications(inputYAML)

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

	// Post-process the minified YAML
	minifiedYAML = postProcessMinifications(minifiedYAML)

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

		// Minify boolean and null values based on string comparison
		if node.Value == "true" {
			node.Value = "y"
		} else if node.Value == "false" {
			node.Value = "n"
		} else if node.Value == "null" {
			node.Value = "~"
		}
	}
}

func preprocessMinifications(input string) string {
	// Regular expressions to identify standalone booleans and nulls
	trueRegex := regexp.MustCompile(`(?m)(^|\s)true($|\s)`)
	falseRegex := regexp.MustCompile(`(?m)(^|\s)false($|\s)`)
	nullRegex := regexp.MustCompile(`(?m)(^|\s)null($|\s)`)

	// Replace true, false, and null with 'y', 'n', and '~'
	input = trueRegex.ReplaceAllStringFunc(input, func(m string) string {
		return strings.Replace(m, "true", "y", -1)
	})
	input = falseRegex.ReplaceAllStringFunc(input, func(m string) string {
		return strings.Replace(m, "false", "n", -1)
	})
	input = nullRegex.ReplaceAllStringFunc(input, func(m string) string {
		return strings.Replace(m, "null", "~", -1)
	})

	return input
}

func postProcessMinifications(input string) string {
	var lines []string
	inMapping := false

	for _, line := range strings.Split(input, "\n") {
		trimmedLine := strings.TrimSpace(line)

		// Detect if we're in a mapping block
		if strings.HasSuffix(trimmedLine, ":") {
			inMapping = true
		} else if trimmedLine == "" {
			inMapping = false
		}

		// Apply minification only to non-indented key-value pairs and within mappings
		if !strings.HasPrefix(line, "  ") && keyValueRegex.MatchString(trimmedLine) {
			if inMapping && !strings.HasSuffix(trimmedLine, ":") {
				line = keyValueRegex.ReplaceAllString(line, "$1:$2")
			} else if !inMapping {
				line = keyValueRegex.ReplaceAllString(line, "$1:$2")
			}
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

var keyValueRegex = regexp.MustCompile(`^\s*([^:\n]+):\s+(.*)$`)
