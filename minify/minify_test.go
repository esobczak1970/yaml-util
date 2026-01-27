// minify/minify_test.go
package minify

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMinifyBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple key-value",
			input:    "key: value\n",
			expected: "key:value\n",
		},
		{
			name:     "Multiple key-values",
			input:    "first: 1\nsecond: 2\n",
			expected: "first:1\nsecond:2\n",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Whitespace only",
			input:    "   \n\n  \n",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Minify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestMinifyMappings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Small inline mapping",
			input: `mappings:
  key1: value1
  key2: value2
`,
			expected: "mappings: {key1: value1, key2: value2}\n",
		},
		{
			name: "Nested mappings",
			input: `config:
  database:
    host: localhost
    port: 5432
`,
			expected: "config:\n  database: {host: localhost, port: 5432}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Minify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			// Normalize whitespace for comparison
			if normalizeWhitespace(result) != normalizeWhitespace(tc.expected) {
				t.Errorf("Expected:\n%q\nGot:\n%q", tc.expected, result)
			}
		})
	}
}

func TestMinifySequences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Short sequence",
			input: `items:
  - item1
  - item2
  - item3
`,
			expected: "items: [item1, item2, item3]\n",
		},
		{
			name: "Long sequence stays expanded",
			input: `items:
  - item1
  - item2
  - item3
  - item4
  - item5
`,
			expected: "items:\n  - item1\n  - item2\n  - item3\n  - item4\n  - item5\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Minify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if normalizeWhitespace(result) != normalizeWhitespace(tc.expected) {
				t.Errorf("Expected:\n%q\nGot:\n%q", tc.expected, result)
			}
		})
	}
}

func TestMinifyBooleanAndNull(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "True to y",
			input:    "enabled: true\n",
			expected: "enabled:y\n",
		},
		{
			name:     "False to n",
			input:    "enabled: false\n",
			expected: "enabled:n\n",
		},
		{
			name:     "Null to tilde",
			input:    "value: null\n",
			expected: "value:~\n",
		},
		{
			name:     "Multiple booleans",
			input:    "flag1: true\nflag2: false\nflag3: null\n",
			expected: "flag1:y\nflag2:n\nflag3:~\n",
		},
		{
			name:     "Boolean in string preserved",
			input:    "text: \"true value\"\n",
			expected: "text:\"true value\"\n",
		},
		{
			name:     "Boolean as part of word",
			input:    "keyword: truevalue\n",
			expected: "keyword:truevalue\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Minify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestMinifyQuotedStrings(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Quoted string with spaces",
			input:   "text: \"value with spaces\"\n",
			wantErr: false,
		},
		{
			name:    "Quoted string with special chars",
			input:   "text: \"special: @#%\"\n",
			wantErr: false,
		},
		{
			name:    "Single quoted string",
			input:   "text: 'single quotes'\n",
			wantErr: false,
		},
		{
			name:    "String with colon",
			input:   "sentence: \"Time: an illusion\"\n",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Minify(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tc.wantErr, err)
			}
			if !tc.wantErr && result == "" {
				t.Error("Expected non-empty result")
			}
		})
	}
}

func TestMinifyAnchorsAndAliases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "Simple anchor and alias",
			input: `defaults: &defaults
  key1: value1
  key2: value2

config:
  <<: *defaults
  key3: value3
`,
			contains: []string{"&defaults", "*defaults"},
		},
		{
			name: "Multiple aliases",
			input: `base: &base
  name: test

ref1:
  <<: *base
  
ref2:
  <<: *base
`,
			contains: []string{"&base", "*base"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Minify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			for _, substr := range tc.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected result to contain %q, got:\n%s", substr, result)
				}
			}
		})
	}
}

func TestMinifyComplexStructures(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "Deeply nested structure",
			input: `root:
  level1:
    level2:
      level3:
        key: value
`,
			wantErr: false,
		},
		{
			name: "Mixed sequences and mappings",
			input: `servers:
  - name: web1
    host: 192.168.1.1
    port: 80
  - name: web2
    host: 192.168.1.2
    port: 80
`,
			wantErr: false,
		},
		{
			name: "Complex real-world config",
			input: `app:
  name: myapp
  version: 1.0.0
  database:
    host: localhost
    port: 5432
    name: mydb
  features:
    - auth
    - api
    - frontend
`,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Minify(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tc.wantErr, err)
			}
			if !tc.wantErr {
				// Verify we can parse the result back
				var testNode interface{}
				if err := yaml.Unmarshal([]byte(result), &testNode); err != nil {
					t.Errorf("Result is not valid YAML: %v\nResult:\n%s", err, result)
				}
			}
		})
	}
}

func TestMinifyErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Invalid YAML - unclosed quote",
			input:   "key: \"value\n",
			wantErr: true,
		},
		{
			name:    "Invalid YAML - bad indentation",
			input:   "key:\n value\n  nested",
			wantErr: true,
		},
		{
			name:    "Invalid YAML - tabs mixed with spaces",
			input:   "key:\n\tvalue",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Minify(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("wantErr=%v, got err=%v", tc.wantErr, err)
			}
		})
	}
}

func TestMinifyWhitespaceHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Leading whitespace",
			input:    "   key: value\n",
			expected: "key:value\n",
		},
		{
			name:     "Trailing whitespace",
			input:    "key: value   \n",
			expected: "key:value\n",
		},
		{
			name:     "Multiple blank lines",
			input:    "key1: value1\n\n\nkey2: value2\n",
			expected: "key1:value1\nkey2:value2\n",
		},
		{
			name:     "Mixed whitespace",
			input:    "  key1: value1  \n\n  key2: value2  \n",
			expected: "key1:value1\nkey2:value2\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Minify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestMinifyPreservesValidity(t *testing.T) {
	inputs := []string{
		"simple: value\n",
		"list:\n  - a\n  - b\n  - c\n",
		"nested:\n  key: value\n",
		"mixed:\n  - item: value\n  - item: value\n",
		"anchor: &ref\n  key: value\nalias:\n  <<: *ref\n",
	}

	for i, input := range inputs {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			// Parse original
			var original interface{}
			if err := yaml.Unmarshal([]byte(input), &original); err != nil {
				t.Fatalf("Original YAML is invalid: %v", err)
			}

			// Minify
			minified, err := Minify(input)
			if err != nil {
				t.Fatalf("Minify failed: %v", err)
			}

			// Parse minified
			var result interface{}
			if err := yaml.Unmarshal([]byte(minified), &result); err != nil {
				t.Errorf("Minified YAML is invalid: %v\nMinified:\n%s", err, minified)
			}
		})
	}
}

// Helper function to normalize whitespace for flexible comparison
func normalizeWhitespace(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, "\n")
}
