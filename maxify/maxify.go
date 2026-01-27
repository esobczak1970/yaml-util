// /maxify/maxify.go
package maxify

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// Maxify expands YAML into a human-readable format.
func Maxify(inputYAML string) (string, error) {
	var rootNode yaml.Node

	// Unmarshal the compact YAML into a structured format.
	err := yaml.Unmarshal([]byte(inputYAML), &rootNode)
	if err != nil {
		return "", err
	}

	// Ensure mappings and lists are properly expanded.
	expandNode(&rootNode)

	// Convert back to YAML with proper formatting.
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(4) // Indent for readability
	defer encoder.Close()

	if err := encoder.Encode(&rootNode); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// expandNode recursively ensures proper formatting for mappings and sequences.
func expandNode(node *yaml.Node) {
	switch node.Kind {
	case yaml.MappingNode:
		node.Style = 0 // Ensures multi-line formatting
		for i := 0; i < len(node.Content); i++ {
			expandNode(node.Content[i])
		}
	case yaml.SequenceNode:
		node.Style = 0 // Expands lists properly
		for _, item := range node.Content {
			expandNode(item)
		}
	case yaml.ScalarNode:
		node.Style = 0 // Ensures correct formatting for values
	}
}

// formatNode recursively ensures proper formatting of mappings and sequences.
func formatNode(node *yaml.Node) {
	switch node.Kind {
	case yaml.MappingNode:
		// Expand all inline mappings into multi-line key-value pairs
		for i := 0; i < len(node.Content); i++ {
			formatNode(node.Content[i])
		}
	case yaml.SequenceNode:
		// Ensure sequences are formatted properly
		for _, item := range node.Content {
			formatNode(item)
		}
	case yaml.ScalarNode:
		// Convert scalars to plain formatting (fixes inline issues)
		node.Style = yaml.LiteralStyle
	}
}
