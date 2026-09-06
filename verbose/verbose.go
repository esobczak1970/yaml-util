// verbose/verbose.go
package verbose

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultIndent = 4

// VerboseOptions controls comment generation behavior.
type VerboseOptions struct {
	// AddTypeComments adds comments describing the type of each element.
	AddTypeComments bool
	// AddStructureComments adds comments about the structure (mapping/sequence).
	AddStructureComments bool
	// AddExamples adds example comments for empty or default values.
	AddExamples bool
	// Indent sets the indentation level.
	Indent int
}

// DefaultOptions returns the default verbose options.
func DefaultOptions() VerboseOptions {
	return VerboseOptions{
		AddExamples:          false,
		AddStructureComments: true,
		AddTypeComments:      true,
		Indent:               defaultIndent,
	}
}

// MakeVerbose adds comments and structure hints to YAML for better understanding.
func MakeVerbose(inputYAML string) (string, error) {
	return MakeVerboseWithOptions(inputYAML, DefaultOptions())
}

// MakeVerboseWithOptions adds comments with custom options.
func MakeVerboseWithOptions(inputYAML string, opts VerboseOptions) (string, error) {
	if strings.TrimSpace(inputYAML) == "" {
		return "", nil
	}

	if err := validateIndentation(inputYAML); err != nil {
		return "", err
	}

	var rootNode yaml.Node
	if err := yaml.Unmarshal([]byte(inputYAML), &rootNode); err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := processNode(&buffer, &rootNode, 0, "", opts); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func processNode(buf *bytes.Buffer, node *yaml.Node, depth int, path string, opts VerboseOptions) error {
	if node == nil {
		return nil
	}

	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			return nil
		}
		if opts.AddStructureComments {
			buf.WriteString("# YAML Document\n")
		}
		return processNode(buf, node.Content[0], depth, path, opts)
	case yaml.MappingNode:
		return processMappingNode(buf, node, depth, path, opts)
	case yaml.SequenceNode:
		return processSequenceNode(buf, node, depth, path, opts)
	case yaml.ScalarNode:
		// path is not needed for scalars, so don't pass it down.
		return processScalarNode(buf, node, depth, opts)
	case yaml.AliasNode:
		indent := strings.Repeat(" ", depth*opts.Indent)
		buf.WriteString(fmt.Sprintf("%s%s", indent, formatAliasValue(node)))
		if opts.AddTypeComments {
			buf.WriteString(" # Alias reference")
		}
		buf.WriteString("\n")
	default:
	}
	return nil
}

func processMappingNode(buf *bytes.Buffer, node *yaml.Node, depth int, path string, opts VerboseOptions) error {
	indent := strings.Repeat(" ", depth*opts.Indent)

	if opts.AddStructureComments {
		buf.WriteString(fmt.Sprintf("%s# Mapping\n", indent))
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		key := keyNode.Value
		currentPath := path
		if currentPath != "" {
			currentPath += "."
		}
		currentPath += key

		buf.WriteString(fmt.Sprintf("%s%s:", indent, key))

		if valueNode == nil {
			continue
		}
		if valueNode.Anchor != "" {
			buf.WriteString(fmt.Sprintf(" &%s", valueNode.Anchor))
		}

		switch valueNode.Kind {
		case yaml.ScalarNode:
			buf.WriteString(" ")
			buf.WriteString(formatScalarValue(valueNode))

			if opts.AddTypeComments {
				comment := getScalarTypeComment(valueNode, opts)
				if comment != "" {
					buf.WriteString(fmt.Sprintf(" # %s", comment))
				}
				if valueNode.Anchor != "" {
					buf.WriteString(fmt.Sprintf(" # anchor: %s", valueNode.Anchor))
				}
			}
			buf.WriteString("\n")

		case yaml.MappingNode:
			if opts.AddTypeComments {
				buf.WriteString(fmt.Sprintf(" # Nested mapping (%d keys)", len(valueNode.Content)/2))
				if valueNode.Anchor != "" {
					buf.WriteString(fmt.Sprintf(" # anchor: %s", valueNode.Anchor))
				}
			}
			buf.WriteString("\n")
			if err := processMappingNode(buf, valueNode, depth+1, currentPath, opts); err != nil {
				return err
			}

		case yaml.SequenceNode:
			if opts.AddTypeComments {
				buf.WriteString(fmt.Sprintf(" # List with %d items", len(valueNode.Content)))
				if valueNode.Anchor != "" {
					buf.WriteString(fmt.Sprintf(" # anchor: %s", valueNode.Anchor))
				}
			}
			buf.WriteString("\n")
			if err := processSequenceNode(buf, valueNode, depth+1, currentPath, opts); err != nil {
				return err
			}

		case yaml.AliasNode:
			buf.WriteString(fmt.Sprintf(" %s", formatAliasValue(valueNode)))
			if opts.AddTypeComments {
				buf.WriteString(" # Alias reference")
			}
			buf.WriteString("\n")

		default:
			buf.WriteString("\n")
		}
	}

	return nil
}

func processSequenceNode(buf *bytes.Buffer, node *yaml.Node, depth int, path string, opts VerboseOptions) error {
	indent := strings.Repeat(" ", depth*opts.Indent)

	if opts.AddStructureComments {
		buf.WriteString(fmt.Sprintf("%s# Sequence\n", indent))
	}

	for idx, item := range node.Content {
		currentPath := fmt.Sprintf("%s[%d]", path, idx)

		buf.WriteString(fmt.Sprintf("%s-", indent))

		if item == nil {
			continue
		}
		if item.Anchor != "" {
			buf.WriteString(fmt.Sprintf(" &%s", item.Anchor))
		}

		switch item.Kind {
		case yaml.ScalarNode:
			buf.WriteString(" ")
			buf.WriteString(formatScalarValue(item))

			if opts.AddTypeComments {
				comment := getScalarTypeComment(item, opts)
				if comment != "" {
					buf.WriteString(fmt.Sprintf(" # %s", comment))
				}
				if item.Anchor != "" {
					buf.WriteString(fmt.Sprintf(" # anchor: %s", item.Anchor))
				}
			}
			buf.WriteString("\n")

		case yaml.MappingNode:
			if opts.AddTypeComments {
				buf.WriteString(fmt.Sprintf(" # Item %d: mapping (%d keys)", idx+1, len(item.Content)/2))
				if item.Anchor != "" {
					buf.WriteString(fmt.Sprintf(" # anchor: %s", item.Anchor))
				}
			}
			buf.WriteString("\n")
			if err := processMappingNode(buf, item, depth+1, currentPath, opts); err != nil {
				return err
			}

		case yaml.SequenceNode:
			if opts.AddTypeComments {
				buf.WriteString(fmt.Sprintf(" # Nested list (%d items)", len(item.Content)))
				if item.Anchor != "" {
					buf.WriteString(fmt.Sprintf(" # anchor: %s", item.Anchor))
				}
			}
			buf.WriteString("\n")
			if err := processSequenceNode(buf, item, depth+1, currentPath, opts); err != nil {
				return err
			}

		case yaml.AliasNode:
			buf.WriteString(fmt.Sprintf(" %s", formatAliasValue(item)))
			if opts.AddTypeComments {
				buf.WriteString(" # Alias reference")
			}
			buf.WriteString("\n")

		default:
			buf.WriteString("\n")
		}
	}

	return nil
}

func processScalarNode(buf *bytes.Buffer, node *yaml.Node, depth int, opts VerboseOptions) error {
	indent := strings.Repeat(" ", depth*opts.Indent)
	buf.WriteString(fmt.Sprintf("%s%s", indent, formatScalarValue(node)))

	if opts.AddTypeComments {
		comment := getScalarTypeComment(node, opts)
		if comment != "" {
			buf.WriteString(fmt.Sprintf(" # %s", comment))
		}
		if node.Anchor != "" {
			buf.WriteString(fmt.Sprintf(" # anchor: %s", node.Anchor))
		}
	}

	buf.WriteString("\n")
	return nil
}

func formatAliasValue(node *yaml.Node) string {
	if node == nil || node.Kind != yaml.AliasNode {
		return "*"
	}
	if node.Alias != nil && node.Alias.Anchor != "" {
		return "*" + node.Alias.Anchor
	}
	if node.Value != "" {
		return "*" + node.Value
	}
	return "*"
}

func formatScalarValue(node *yaml.Node) string {
	value := node.Value

	anchor := ""
	if node.Anchor != "" {
		anchor = fmt.Sprintf("&%s ", node.Anchor)
	}

	switch node.Style {
	case yaml.DoubleQuotedStyle:
		return fmt.Sprintf("%s\"%s\"", anchor, value)
	case yaml.SingleQuotedStyle:
		return fmt.Sprintf("%s'%s'", anchor, value)
	default:
		if needsQuotes(value) {
			return fmt.Sprintf("%s\"%s\"", anchor, value)
		}
		return anchor + value
	}
}

func needsQuotes(value string) bool {
	if value == "" {
		return true
	}
	if value != strings.TrimSpace(value) {
		return true
	}
	return strings.ContainsAny(value, ":{}[]|>*&!%@`#")
}

func getScalarTypeComment(node *yaml.Node, opts VerboseOptions) string {
	value := node.Value

	switch {
	case value == "true" || value == "false" || value == "y" || value == "n":
		return "boolean"
	case value == "null" || value == "~":
		return "null value"
	case isInteger(value):
		return "integer"
	case isFloat(value):
		return "float"
	default:
		if len(value) > 50 {
			return "string (long)"
		}
		if value == "" && opts.AddExamples {
			return "string (empty)"
		}
		return "string"
	}
}

func isInteger(s string) bool {
	if s == "" {
		return false
	}

	start := 0
	if s[0] == '-' || s[0] == '+' {
		if len(s) == 1 {
			return false
		}
		start = 1
	}

	for i := start; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}

	return true
}

func isFloat(s string) bool {
	if s == "" {
		return false
	}

	start := 0
	if s[0] == '-' || s[0] == '+' {
		if len(s) == 1 {
			return false
		}
		start = 1
	}

	dotFound := false
	for i := start; i < len(s); i++ {
		if s[i] == '.' {
			if dotFound {
				return false
			}
			dotFound = true
			continue
		}
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}

	return dotFound
}

func validateIndentation(input string) error {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.ContainsRune(line, '\t') {
			return errors.New("invalid indentation: tabs are not allowed")
		}
		leading := countLeadingSpaces(line)
		if leading%2 != 0 {
			return errors.New("invalid indentation: indentation must be a multiple of 2 spaces")
		}
	}
	return nil
}

func countLeadingSpaces(s string) int {
	n := 0
	for n < len(s) && s[n] == ' ' {
		n++
	}
	return n
}
