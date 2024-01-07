package minify

import (
	"testing"
)

// TestBasicYAMLStructure ensures basic YAML structures are preserved.
func TestBasicYAMLStructure(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"Simple Key-Value", "key: value\n", "key: value\n"},
		{"Multiple Key-Values", "first: 1\nsecond: 2\n", "first: 1\nsecond: 2\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			minified, err := Minify(tc.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if minified != tc.output {
				t.Errorf("Expected %q, got %q", tc.output, minified)
			}
		})
	}
}

// TestWhitespaceAndNewlineMinification checks if unnecessary whitespaces and newlines are removed.
func TestWhitespaceAndNewlineMinification(t *testing.T) {
	input := "   key: value   \n\n"
	expectedOutput := "key: value\n"

	minified, err := Minify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if minified != expectedOutput {
		t.Errorf("Expected %q, got %q", expectedOutput, minified)
	}
}

// TestErrorHandling ensures that invalid YAML results in an error.
func TestErrorHandling(t *testing.T) {
	invalidYAML := "key: value\ninvalid_yaml"

	_, err := Minify(invalidYAML)
	if err == nil {
		t.Errorf("Expected an error for invalid YAML, but got none")
	}
}

// TestMinificationOfSequencesAndMappings checks minification of sequences and mappings.
func TestMinificationOfSequencesAndMappings(t *testing.T) {
	input := `
sequences:
  - item1
  - item2
mappings:
  key1: value1
  key2: value2
`
	expectedOutput := "sequences:\n  - item1\n  - item2\nmappings:\n  key1: value1\n  key2: value2\n"
	minified, err := Minify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if minified != expectedOutput {
		t.Errorf("Expected %q, got %q", expectedOutput, minified)
	}
}

// TestWhitespaceAndNewlineRemoval checks if unnecessary whitespaces and newlines are properly removed.
func TestWhitespaceAndNewlineRemoval(t *testing.T) {
	input := `
    key:
        - value1
        - value2
`
	// Correctly minified YAML should maintain the necessary indentation
	expectedOutput := "key:\n  - value1\n  - value2\n"

	minified, err := Minify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if minified != expectedOutput {
		t.Errorf("Expected %q, got %q", expectedOutput, minified)
	}
}
