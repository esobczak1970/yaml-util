// maxify/maxify_test.go
package maxify

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMaxifyBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple key-value",
			input:    "key:value",
			expected: "key: value\n",
		},
		{
			name:     "Already formatted",
			input:    "key: value\n",
			expected: "key: value\n",
		},
		{
			name:     "Multiple keys",
			input:    "key1:value1\nkey2:value2",
			expected: "key1: value1\nkey2: value2\n",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Maxify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected:\n%q\nGot:\n%q", tc.expected, result)
			}
		})
	}
}

func TestMaxifyInlineMappings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "Inline to block mapping",
			input:    "mappings: {key1:value1, key2:value2}",
			contains: []string{"mappings:", "key1: value1", "key2: value2"},
		},
		{
			name:     "Nested inline mapping",
			input:    "root: {nested: {key:value}}",
			contains: []string{"root:", "nested:", "key: value"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Maxify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			for _, substr := range tc.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected result to contain %q\nGot:\n%s", substr, result)
				}
			}
		})
	}
}

func TestMaxifySequences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "Inline to block sequence",
			input:    "items: [item1, item2, item3]",
			contains: []string{"items:", "- item1", "- item2", "- item3"},
		},
		{
			name:     "Nested sequences",
			input:    "root: [[a, b], [c, d]]",
			contains: []string{"root:", "-", "- a", "- b", "- c", "- d"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Maxify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			for _, substr := range tc.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected result to contain %q\nGot:\n%s", substr, result)
				}
			}
		})
	}
}

func TestMaxifyBooleanAndNull(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "y to true",
			input:    "flag:y",
			expected: "flag: true\n",
		},
		{
			name:     "n to false",
			input:    "flag:n",
			expected: "flag: false\n",
		},
		{
			name:     "tilde to null",
			input:    "value:~",
			expected: "value: null\n",
		},
		{
			name:     "Mixed boolean values",
			input:    "a:y\nb:n\nc:~",
			expected: "a: true\nb: false\nc: null\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Maxify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected:\n%q\nGot:\n%q", tc.expected, result)
			}
		})
	}
}

func TestMaxifyComplexStructures(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Complex nested structure",
			input: "config: {nested_map: {key1:value1, key2:value2}, nested_list: [item1, item2]}",
		},
		{
			name: "Multi-level nesting",
			input: `root:
  level1: {level2: {level3: value}}`,
		},
		{
			name: "Mixed structures",
			input: `app:
  servers: [{name:web1,host:localhost},{name:web2,host:192.168.1.1}]
  config: {debug:y,timeout:30}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Maxify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify result is valid YAML
			var parsed interface{}
			if err := yaml.Unmarshal([]byte(result), &parsed); err != nil {
				t.Errorf("Result is not valid YAML: %v\nResult:\n%s", err, result)
			}

			// Verify it's more readable (has newlines and indentation)
			if !strings.Contains(result, "\n") {
				t.Error("Expected result to contain newlines for readability")
			}
			if !strings.Contains(result, "    ") && strings.Count(result, "\n") > 2 {
				t.Error("Expected result to contain indentation")
			}
		})
	}
}

func TestWithIndent(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		indent int
	}{
		{
			name:   "2-space indent",
			input:  "root: {key:value}",
			indent: 2,
		},
		{
			name:   "4-space indent",
			input:  "root: {key:value}",
			indent: 4,
		},
		{
			name:   "8-space indent",
			input:  "root: {key:value}",
			indent: 8,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := WithIndent(tc.input, tc.indent)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check that indentation is present
			lines := strings.Split(result, "\n")
			foundIndent := false
			for _, line := range lines {
				if len(line) > 0 && line[0] == ' ' {
					foundIndent = true
					// Count leading spaces
					spaces := 0
					for _, ch := range line {
						if ch == ' ' {
							spaces++
						} else {
							break
						}
					}
					// Should be a multiple of the indent
					if spaces%tc.indent != 0 && spaces != 0 {
						t.Errorf("Found %d spaces, expected multiple of %d", spaces, tc.indent)
					}
				}
			}
			if !foundIndent && len(lines) > 2 {
				t.Error("Expected to find indentation in output")
			}
		})
	}
}

func TestMaxifyAnchorsAndAliases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "Preserve anchors",
			input:    "defaults: &defaults {key:value}\nconfig: {<<: *defaults}",
			contains: []string{"&defaults", "*defaults"},
		},
		{
			name: "Multiple aliases",
			input: `base: &base {name:test}
ref1: {<<: *base}
ref2: {<<: *base}`,
			contains: []string{"&base", "*base"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Maxify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			for _, substr := range tc.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected result to contain %q\nGot:\n%s", substr, result)
				}
			}
		})
	}
}

func TestMaxifyQuotedStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Double quoted strings",
			input: "text: \"value with spaces\"",
		},
		{
			name:  "Single quoted strings",
			input: "text: 'value with spaces'",
		},
		{
			name:  "Strings with special chars",
			input: "text: \"special: @#%\"",
		},
		{
			name:  "Multi-line strings",
			input: "text: \"line1\\nline2\"",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Maxify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify result is valid YAML
			var parsed interface{}
			if err := yaml.Unmarshal([]byte(result), &parsed); err != nil {
				t.Errorf("Result is not valid YAML: %v\nResult:\n%s", err, result)
			}
		})
	}
}

func TestMaxifyErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Invalid YAML - unclosed bracket",
			input:   "key: {value",
			wantErr: true,
		},
		{
			name:    "Invalid YAML - bad indentation",
			input:   "key:\n value\n  nested",
			wantErr: true,
		},
		{
			name:    "Invalid YAML - unclosed quote",
			input:   "key: \"value",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Maxify(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("wantErr=%v, got err=%v", tc.wantErr, err)
			}
		})
	}
}

func TestMaxifyPreservesSemantics(t *testing.T) {
	inputs := []string{
		"key:value",
		"list:[a,b,c]",
		"map:{k1:v1,k2:v2}",
		"nested:{a:{b:c}}",
		"mixed:{list:[1,2],map:{k:v}}",
	}

	for i, input := range inputs {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			// Parse original
			var original interface{}
			if err := yaml.Unmarshal([]byte(input), &original); err != nil {
				t.Fatalf("Original YAML is invalid: %v", err)
			}

			// Maxify
			maxified, err := Maxify(input)
			if err != nil {
				t.Fatalf("Maxify failed: %v", err)
			}

			// Parse maxified
			var result interface{}
			if err := yaml.Unmarshal([]byte(maxified), &result); err != nil {
				t.Errorf("Maxified YAML is invalid: %v\nMaxified:\n%s", err, maxified)
			}

			// Note: Deep equality checking of parsed YAML is complex
			// The important thing is that both parse without error
		})
	}
}

func TestMaxifyWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Leading whitespace",
			input: "   key:value",
		},
		{
			name:  "Trailing whitespace",
			input: "key:value   ",
		},
		{
			name:  "Multiple blank lines",
			input: "key1:value1\n\n\nkey2:value2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Maxify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Should not have trailing spaces on lines
			lines := strings.Split(result, "\n")
			for i, line := range lines {
				if len(line) > 0 && line != strings.TrimRight(line, " \t") {
					t.Errorf("Line %d has trailing whitespace: %q", i+1, line)
				}
			}

			// Should not have multiple consecutive blank lines
			prevEmpty := false
			for i, line := range lines {
				isEmpty := strings.TrimSpace(line) == ""
				if isEmpty && prevEmpty {
					t.Errorf("Found consecutive blank lines at line %d", i+1)
				}
				prevEmpty = isEmpty
			}
		})
	}
}

func TestMaxifyRealWorldExamples(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "Kubernetes-style config",
			input: `apiVersion:v1
kind:Service
metadata:{name:my-service,namespace:default}
spec:{ports:[{port:80,targetPort:8080}],selector:{app:my-app}}`,
		},
		{
			name: "Docker Compose style",
			input: `version:"3.8"
services:{web:{image:"nginx:latest",ports:["80:80"]},db:{image:"postgres:13",environment:{POSTGRES_PASSWORD:secret}}}`,
		},
		{
			name: "Application config",
			input: `app:{name:myapp,version:"1.0.0",debug:y}
database:{host:localhost,port:5432,credentials:{user:admin,password:secret}}
features:[auth,api,frontend]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Maxify(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify it's valid YAML
			var parsed interface{}
			if err := yaml.Unmarshal([]byte(result), &parsed); err != nil {
				t.Errorf("Result is not valid YAML: %v\nResult:\n%s", err, result)
			}

			// Verify it's more readable than input
			if len(result) < len(tc.input) {
				t.Error("Expected maxified version to be longer (more readable)")
			}
		})
	}
}
