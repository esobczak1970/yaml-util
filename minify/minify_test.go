// /minify/minify_test.go
package minify

import (
	"strings"
	"testing"
)

// TestMinifyBasic ensures that minification removes unnecessary spaces.
func TestMinifyBasic(t *testing.T) {
	input := "key: value\n"
	expectedOutput := "key:value\n"

	minified, err := Minify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if strings.TrimSpace(minified) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected %q, got %q", expectedOutput, minified)
	}
}

// TestMinificationOfMappings fixes inline format handling.
func TestMinificationOfMappings(t *testing.T) {
	input := `
mappings:
  key1: value1
  key2: value2
`
	expectedOutput := "mappings:{key1:value1, key2:value2}\n"

	minified, err := Minify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if minified != expectedOutput {
		t.Errorf("Expected %q, got %q", expectedOutput, minified)
	}
}

// TestAnchorAliasMinification ensures anchors & aliases remain valid.
func TestAnchorAliasMinification(t *testing.T) {
	input := `
defaults: &defaults
  key1: value1
  key2: value2

config:
  <<: *defaults
  key3: value3
`
	expectedOutput := "defaults:&defaults {key1:value1, key2:value2}\nconfig:{<<:*defaults, key3:value3}\n"

	minified, err := Minify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if minified != expectedOutput {
		t.Errorf("Expected %q, got %q", expectedOutput, minified)
	}
}

// TestComplexDataStructures ensures nested maps/lists are inline.
func TestComplexDataStructures(t *testing.T) {
	input := `
top-level:
  nested-map:
    key1: value1
    key2: value2
  nested-list:
    - list-item1
    - list-item2
`
	expectedOutput := "top-level:{nested-map:{key1:value1, key2:value2}, nested-list:[list-item1, list-item2]}\n"

	minified, err := Minify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if minified != expectedOutput {
		t.Errorf("Expected %q, got %q", expectedOutput, minified)
	}
}

// TestSafeModePreservesValidity checks safe mode behavior.
func TestSafeModePreservesValidity(t *testing.T) {
	input := `
nested:
  key: "value with spaces"
  another: "special characters: @#%"
`
	expectedOutput := "nested:{key:\"value with spaces\", another:\"special characters: @#%\"}\n"

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

// TestBasicYAMLStructure ensures basic YAML structures are preserved.
func TestBasicYAMLStructure(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"Simple Key-Value", "key: value\n", "key:value\n"},
		{"Multiple Key-Values", "first: 1\nsecond: 2\n", "first:1\nsecond:2\n"},
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
	expectedOutput := "key:value\n"

	minified, err := Minify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if minified != expectedOutput {
		t.Errorf("Expected %q, got %q", expectedOutput, minified)
	}
}

// TestMinificationOfSequences checks minification of sequences.
func TestMinificationOfSequences(t *testing.T) {
	input := `
sequences:
  - item1
  - item2
`
	expectedOutput := "sequences:\n  - item1\n  - item2\n"

	minified, err := Minify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if minified != expectedOutput {
		t.Errorf("Expected %q, got %q", expectedOutput, minified)
	}
}

// TestBooleanAndNullMinification checks boolean and null value minification.
func TestBooleanAndNullMinification(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"True to 'y'", "boolean: true\n", "boolean:y\n"},
		{"False to 'n'", "boolean: false\n", "boolean:n\n"},
		{"Null to '~'", "value: null\n", "value:~\n"},
		{"Preserve '~'", "value: ~\n", "value:~\n"},
		{"Inside quotes", "text: \"true is not y\"\n", "text:\"true is not y\"\n"},
		{"As part of a word", "keyword: truevalue\n", "keyword:truevalue\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			minified, err := Minify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if minified != tc.output {
				t.Errorf("Expected %q, got %q", tc.output, minified)
			}
		})
	}
}

// TestKeyValueMinificationOnOwnLine checks key-value spacing minification.
func TestKeyValueMinificationOnOwnLine(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"Single Key-Value Pair", "key: value\n", "key:value\n"},
		{"Multiple Key-Value Pairs", "first: 1\nsecond: 2\n", "first:1\nsecond:2\n"},
		{"Colon in Middle of String", "sentence: \"Time: an illusion\"\n", "sentence:\"Time: an illusion\"\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			minified, err := Minify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if minified != tc.output {
				t.Errorf("Expected %q, got %q", tc.output, minified)
			}
		})
	}
}
