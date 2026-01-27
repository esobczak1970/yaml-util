// minify/minify.go
package minify

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

// Minify reduces YAML to its most compact valid form.
func Minify(inputYAML string) (string, error) {
	var rootNode yaml.Node

	// Parse YAML structure
	err := yaml.Unmarshal([]byte(inputYAML), &rootNode)
	if err != nil {
		return "", err
	}

	// Apply minification transformations
	minifyNode(&rootNode)

	// Convert back to YAML with minimal indentation
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(0) // Ensures tight formatting
	defer encoder.Close()

	if err := encoder.Encode(&rootNode); err != nil {
		return "", err
	}

	// Trim unnecessary spaces from output
	return cleanupFormatting(buffer.String()), nil
}

// cleanupFormatting removes unnecessary spaces and optimizes key-value pairs.
func cleanupFormatting(input string) string {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		// Ensure proper spacing for mappings (key:value)
		if strings.Contains(line, ": ") && !strings.Contains(line, "\"") && !strings.Contains(line, "'") {
			lines[i] = strings.Replace(line, ": ", ":", 1)
		}
	}

	// Convert `true/false/null` to `y/n/~` but avoid modifying words like "truevalue"
	output := strings.Join(lines, "\n")
	output = strings.ReplaceAll(output, "\ntrue\n", "\ny\n")
	output = strings.ReplaceAll(output, "\nfalse\n", "\nn\n")
	output = strings.ReplaceAll(output, "\nnull\n", "\n~\n")

	return output
}

// minifyNode applies compact transformations to YAML nodes.
func minifyNode(node *yaml.Node) {
	switch node.Kind {
	case yaml.MappingNode:
		inlineMapping(node)
		for _, n := range node.Content {
			minifyNode(n)
		}
	case yaml.SequenceNode:
		inlineSequence(node)
	case yaml.ScalarNode:
		minifyScalar(node)
	}
}

// inlineMapping converts mappings to inline `{key:value, key2:value2}` format.
func inlineMapping(node *yaml.Node) {
	if node.Kind != yaml.MappingNode || len(node.Content) > 4 {
		return // Avoid excessive inline minifications.
	}

	inlineFormat := make([]string, 0, len(node.Content)/2)
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		inlineFormat = append(inlineFormat, key+":"+node.Content[i+1].Value)
	}

	node.Value = "{" + strings.Join(inlineFormat, ", ") + "}"
	node.Style = yaml.FlowStyle
	node.Content = nil // Clear children since it's now inline
}

// inlineSequence converts short lists to `[item1, item2]` format.
func inlineSequence(node *yaml.Node) {
	if len(node.Content) > 3 { // Keep inline only for short lists.
		return
	}
	values := make([]string, len(node.Content))
	for i, item := range node.Content {
		values[i] = item.Value
	}
	node.Value = "[" + strings.Join(values, ", ") + "]"
	node.Style = yaml.FlowStyle
	node.Content = nil // Clear children since it's now inline
}

// minifyScalar removes unnecessary quotes and normalizes scalar values.
func minifyScalar(node *yaml.Node) {
	node.Value = strings.TrimSpace(node.Value)

	// Remove unnecessary quotes unless needed
	if node.Style == yaml.DoubleQuotedStyle || node.Style == yaml.SingleQuotedStyle {
		node.Style = yaml.FlowStyle
	}

	// Minify boolean and null values safely
	if node.Value == "true" {
		node.Value = "y"
	} else if node.Value == "false" {
		node.Value = "n"
	} else if node.Value == "null" {
		node.Value = "~"
	}
}
