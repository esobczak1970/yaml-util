// Package maxify provides utilities for expanding YAML.
package maxify

import (
	"bytes"
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultIndent = 4

// Maxify expands minified YAML-like input into a human-readable YAML format.
func Maxify(inputYAML string) (string, error) {
	return WithIndent(inputYAML, defaultIndent)
}

// WithIndent expands YAML with a custom indentation level.
func WithIndent(inputYAML string, indent int) (string, error) {
	if strings.TrimSpace(inputYAML) == "" {
		return "", nil
	}

	if err := validateIndentation(inputYAML); err != nil {
		return "", err
	}

	normalized := insertSpaceAfterColon(inputYAML)

	var rootNode yaml.Node
	if err := yaml.Unmarshal([]byte(normalized), &rootNode); err != nil {
		return "", err
	}

	expandNode(&rootNode)

	// Encode the document content (not the DocumentNode) to avoid leading "---".
	nodeToEncode := &rootNode
	if rootNode.Kind == yaml.DocumentNode && len(rootNode.Content) == 1 {
		nodeToEncode = rootNode.Content[0]
	}

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(indent)
	defer encoder.Close()

	if err := encoder.Encode(nodeToEncode); err != nil {
		return "", err
	}

	return ensureTrailingNewline(trimTrailingWhitespace(buffer.String())), nil
}

func expandNode(node *yaml.Node) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.DocumentNode:
		for _, n := range node.Content {
			expandNode(n)
		}
	case yaml.MappingNode, yaml.SequenceNode:
		// Force block style for readability.
		node.Style &^= yaml.FlowStyle
		for _, n := range node.Content {
			expandNode(n)
		}
	case yaml.ScalarNode:
		expandScalar(node)
	case yaml.AliasNode:
		// Preserve aliases as-is.
	}
}

func expandScalar(node *yaml.Node) {
	if node == nil || node.Kind != yaml.ScalarNode {
		return
	}

	// Expand shorthand bool/null in plain scalars.
	if node.Style == 0 {
		switch node.Value {
		case "y", "Y", "yes", "Yes", "YES":
			node.Tag = "!!bool"
			node.Value = "true"
		case "n", "N", "no", "No", "NO":
			node.Tag = "!!bool"
			node.Value = "false"
		case "~":
			node.Tag = "!!null"
			node.Value = "null"
		}
	}

	// Be robust if the parser already tagged it but kept a short form.
	switch node.Tag {
	case "!!bool":
		if node.Value == "y" || node.Value == "Y" {
			node.Value = "true"
		}
		if node.Value == "n" || node.Value == "N" {
			node.Value = "false"
		}
	case "!!null":
		if node.Value == "~" {
			node.Value = "null"
		}
	}
}

// validateIndentation is intentionally permissive:
// - Leading whitespace is allowed.
// - Tabs are rejected.
// - Indentation may increase only after a line that starts a block (ends with ":" or is "-").
// This catches the test case: "key:\n value\n  nested" (nested under a scalar line).
func validateIndentation(input string) error {
	lines := strings.Split(input, "\n")

	prevLine := ""
	prevIndent := 0
	havePrev := false

	for _, raw := range lines {
		if strings.TrimSpace(raw) == "" {
			continue
		}
		if strings.ContainsRune(raw, '\t') {
			return errors.New("invalid indentation: tabs are not allowed")
		}

		indent := countLeadingSpaces(raw)
		line := strings.TrimSpace(raw)

		if havePrev && indent > prevIndent {
			prevTrim := strings.TrimSpace(prevLine)
			if !strings.HasSuffix(prevTrim, ":") && prevTrim != "-" {
				return errors.New("invalid indentation: indentation increased under a non-block line")
			}
		}

		prevLine = line
		prevIndent = indent
		havePrev = true
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

// insertSpaceAfterColon converts "key:value" into "key: value" when the colon is
// clearly acting as a mapping separator (outside of quotes).
func insertSpaceAfterColon(input string) string {
	var b strings.Builder
	b.Grow(len(input))

	inSingle := false
	inDouble := false

	for i := 0; i < len(input); i++ {
		c := input[i]

		if inSingle {
			if c == '\'' {
				inSingle = false
			}
			b.WriteByte(c)
			continue
		}

		if inDouble {
			if c == '\\' && i+1 < len(input) {
				b.WriteByte(c)
				i++
				b.WriteByte(input[i])
				continue
			}
			if c == '"' {
				inDouble = false
			}
			b.WriteByte(c)
			continue
		}

		switch c {
		case '\'':
			inSingle = true
			b.WriteByte(c)
			continue
		case '"':
			inDouble = true
			b.WriteByte(c)
			continue
		case ':':
			b.WriteByte(c)
			if i+1 >= len(input) {
				continue
			}
			next := input[i+1]
			if next == ' ' || next == '\n' || next == '\r' || next == '\t' || next == ',' || next == '}' || next == ']' {
				continue
			}
			b.WriteByte(' ')
			continue
		default:
			b.WriteByte(c)
		}
	}

	return b.String()
}

func trimTrailingWhitespace(input string) string {
	lines := strings.Split(input, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	return strings.Join(lines, "\n")
}

func ensureTrailingNewline(s string) string {
	if s == "" || strings.HasSuffix(s, "\n") {
		return s
	}
	return s + "\n"
}
