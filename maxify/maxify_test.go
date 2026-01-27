// /maxify/maxify_test.go
package maxify

import (
	"strings"
	"testing"
)

// TestMaxifyBasic ensures basic YAML expands properly.
func TestMaxifyBasic(t *testing.T) {
	input := "key:value\n"
	expectedOutput := "key: value\n"

	maxified, err := Maxify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if strings.TrimSpace(maxified) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected %q, got %q", expectedOutput, maxified)
	}
}

// TestMaxifyInlineMappings ensures inline mappings expand correctly.
func TestMaxifyInlineMappings(t *testing.T) {
	input := "mappings: {key1:value1, key2:value2}\n"
	expectedOutput := `mappings:
    key1: value1
    key2: value2
`

	maxified, err := Maxify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if strings.TrimSpace(maxified) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected:\n%q\nGot:\n%q", expectedOutput, maxified)
	}
}

// TestMaxifyLists ensures lists are properly formatted.
func TestMaxifyLists(t *testing.T) {
	input := "items: [item1, item2, item3]\n"
	expectedOutput := `items:
    - item1
    - item2
    - item3
`

	maxified, err := Maxify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if strings.TrimSpace(maxified) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected:\n%q\nGot:\n%q", expectedOutput, maxified)
	}
}

// TestMaxifyComplexStructures ensures nested structures expand correctly.
func TestMaxifyComplexStructures(t *testing.T) {
	input := `
config: {nested_map: {key1:value1, key2:value2}, nested_list: [item1, item2]}
`
	expectedOutput := `config:
    nested_map:
        key1: value1
        key2: value2
    nested_list:
        - item1
        - item2
`

	maxified, err := Maxify(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if strings.TrimSpace(maxified) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected:\n%q\nGot:\n%q", expectedOutput, maxified)
	}
}
