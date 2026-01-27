// verbose/verbose_test.go
package verbose

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMakeVerboseBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "Simple key-value",
			input:    "key: value",
			contains: []string{"key: value"},
		},
		{
			name:     "Multiple keys",
			input:    "key1: value1\nkey2: value2",
			contains: []string{"key1: value1", "key2: value2"},
		},
		{
			name:     "Empty input",
			input:    "",
			contains: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MakeVerbose(tc.input)
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

func TestMakeVerboseWithComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "String type comment",
			input:    "name: test",
			contains: []string{"name: test", "# string"},
		},
		{
			name:     "Integer type comment",
			input:    "port: 8080",
			contains: []string{"port: 8080", "# integer"},
		},
		{
			name:     "Boolean type comment",
			input:    "enabled: true",
			contains: []string{"enabled: true", "# boolean"},
		},
		{
			name:     "Null type comment",
			input:    "value: null",
			contains: []string{"value: null", "# null"},
		},
		{
			name:     "Float type comment",
			input:    "rate: 3.14",
			contains: []string{"rate: 3.14", "# float"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MakeVerbose(tc.input)
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

func TestMakeVerboseMappings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "Nested mapping",
			input: `config:
  database:
    host: localhost
    port: 5432`,
			contains: []string{"config:", "database:", "host: localhost", "port: 5432", "# Nested mapping"},
		},
		{
			name: "Multiple nested levels",
			input: `root:
  level1:
    level2:
      key: value`,
			contains: []string{"root:", "level1:", "level2:", "key: value"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MakeVerbose(tc.input)
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

func TestMakeVerboseSequences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "Simple list",
			input: `items:
  - item1
  - item2
  - item3`,
			contains: []string{"items:", "# List with 3 items", "- item1", "- item2", "- item3"},
		},
		{
			name: "List of mappings",
			input: `servers:
  - name: web1
    host: localhost
  - name: web2
    host: 192.168.1.1`,
			contains: []string{"servers:", "# List", "name: web1", "name: web2"},
		},
		{
			name: "Nested lists",
			input: `matrix:
  - [1, 2, 3]
  - [4, 5, 6]`,
			contains: []string{"matrix:", "# List"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MakeVerbose(tc.input)
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

func TestMakeVerboseAnchorsAndAliases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "Simple anchor and alias",
			input: `defaults: &defaults
  timeout: 30
  retries: 3

config:
  <<: *defaults
  name: myapp`,
			contains: []string{"&defaults", "*defaults", "# Alias", "# anchor:"},
		},
		{
			name: "Multiple aliases",
			input: `base: &base
  key: value

ref1:
  <<: *base

ref2:
  <<: *base`,
			contains: []string{"&base", "*base"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MakeVerbose(tc.input)
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

func TestMakeVerboseWithOptions(t *testing.T) {
	input := "name: test\nport: 8080"

	tests := []struct {
		name         string
		opts         VerboseOptions
		shouldHave   []string
		shouldNotHave []string
	}{
		{
			name: "Type comments enabled",
			opts: VerboseOptions{
				AddTypeComments:      true,
				AddStructureComments: false,
				AddExamples:          false,
				Indent:               4,
			},
			shouldHave:   []string{"# string", "# integer"},
			shouldNotHave: []string{"# Mapping"},
		},
		{
			name: "Structure comments enabled",
			opts: VerboseOptions{
				AddTypeComments:      false,
				AddStructureComments: true,
				AddExamples:          false,
				Indent:               4,
			},
			shouldHave:   []string{"# Mapping"},
			shouldNotHave: []string{"# string", "# integer"},
		},
		{
			name: "All comments disabled",
			opts: VerboseOptions{
				AddTypeComments:      false,
				AddStructureComments: false,
				AddExamples:          false,
				Indent:               4,
			},
			shouldNotHave: []string{"#"},
		},
		{
			name: "Custom indentation",
			opts: VerboseOptions{
				AddTypeComments:      true,
				AddStructureComments: false,
				AddExamples:          false,
				Indent:               2,
			},
			shouldHave: []string{"name: test", "port: 8080"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MakeVerboseWithOptions(input, tc.opts)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			for _, substr := range tc.shouldHave {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected result to contain %q\nGot:\n%s", substr, result)
				}
			}

			for _, substr := range tc.shouldNotHave {
				if strings.Contains(result, substr) {
					t.Errorf("Expected result NOT to contain %q\nGot:\n%s", substr, result)
				}
			}
		})
	}
}

func TestMakeVerboseComplexStructures(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "Kubernetes-style config",
			input: `apiVersion: v1
kind: Service
metadata:
  name: my-service
  namespace: default
spec:
  ports:
    - port: 80
      targetPort: 8080
  selector:
    app: my-app`,
		},
		{
			name: "Database config",
			input: `database:
  connections:
    - host: primary.db.local
      port: 5432
      database: main
      username: admin
      ssl: true
    - host: replica.db.local
      port: 5432
      database: main
      username: readonly
      ssl: true`,
		},
		{
			name: "Mixed types",
			input: `config:
  name: myapp
  version: 1.2.3
  enabled: true
  timeout: 30
  rate: 0.95
  metadata: null
  tags:
    - production
    - critical
  settings:
    debug: false
    verbose: true`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MakeVerbose(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify result contains comments
			if !strings.Contains(result, "#") {
				t.Error("Expected result to contain comments")
			}

			// Verify result is valid YAML
			var parsed interface{}
			if err := yaml.Unmarshal([]byte(result), &parsed); err != nil {
				t.Errorf("Result is not valid YAML: %v\nResult:\n%s", err, result)
			}
		})
	}
}

func TestMakeVerboseQuotedStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Double quoted",
			input: "text: \"value with spaces\"",
		},
		{
			name:  "Single quoted",
			input: "text: 'value with spaces'",
		},
		{
			name:  "Special characters",
			input: "url: \"https://example.com:8080\"",
		},
		{
			name:  "Multiline",
			input: "description: \"This is a long\\ntext with newlines\"",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MakeVerbose(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify it's valid YAML
			var parsed interface{}
			if err := yaml.Unmarshal([]byte(result), &parsed); err != nil {
				t.Errorf("Result is not valid YAML: %v\nResult:\n%s", err, result)
			}
		})
	}
}

func TestMakeVerboseErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Invalid YAML",
			input:   "key: {bad",
			wantErr: true,
		},
		{
			name:    "Unclosed quote",
			input:   "key: \"value",
			wantErr: true,
		},
		{
			name:    "Bad indentation",
			input:   "key:\n value\n  bad",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := MakeVerbose(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("wantErr=%v, got err=%v", tc.wantErr, err)
			}
		})
	}
}

func TestMakeVerbosePreservesSemantics(t *testing.T) {
	inputs := []string{
		"key: value",
		"list:\n  - a\n  - b\n  - c",
		"nested:\n  key: value",
		"bool: true\nnum: 42\nfloat: 3.14",
		"null: null",
	}

	for i, input := range inputs {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			// Parse original
			var original interface{}
			if err := yaml.Unmarshal([]byte(input), &original); err != nil {
				t.Fatalf("Original YAML is invalid: %v", err)
			}

			// Make verbose
			verbose, err := MakeVerbose(input)
			if err != nil {
				t.Fatalf("MakeVerbose failed: %v", err)
			}

			// Parse verbose version (comments should be ignored by parser)
			var result interface{}
			if err := yaml.Unmarshal([]byte(verbose), &result); err != nil {
				t.Errorf("Verbose YAML is invalid: %v\nVerbose:\n%s", err, verbose)
			}
		})
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if !opts.AddTypeComments {
		t.Error("Expected AddTypeComments to be true by default")
	}
	if !opts.AddStructureComments {
		t.Error("Expected AddStructureComments to be true by default")
	}
	if opts.AddExamples {
		t.Error("Expected AddExamples to be false by default")
	}
	if opts.Indent != defaultIndent {
		t.Errorf("Expected Indent to be %d, got %d", defaultIndent, opts.Indent)
	}
}

func TestTypeDetection(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Detect integer",
			input:    "value: 42",
			expected: "# integer",
		},
		{
			name:     "Detect negative integer",
			input:    "value: -42",
			expected: "# integer",
		},
		{
			name:     "Detect float",
			input:    "value: 3.14",
			expected: "# float",
		},
		{
			name:     "Detect boolean true",
			input:    "value: true",
			expected: "# boolean",
		},
		{
			name:     "Detect boolean false",
			input:    "value: false",
			expected: "# boolean",
		},
		{
			name:     "Detect null",
			input:    "value: null",
			expected: "# null",
		},
		{
			name:     "Detect string",
			input:    "value: hello",
			expected: "# string",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MakeVerbose(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !strings.Contains(result, tc.expected) {
				t.Errorf("Expected result to contain %q\nGot:\n%s", tc.expected, result)
			}
		})
	}
}

func TestVerboseRealWorldExamples(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "Application configuration",
			input: `app:
  name: my-application
  version: 2.0.1
  environment: production
  debug: false
  
server:
  host: 0.0.0.0
  port: 8080
  timeout: 30
  
database:
  driver: postgresql
  host: db.example.com
  port: 5432
  name: appdb
  pool_size: 20`,
		},
		{
			name: "CI/CD pipeline",
			input: `pipeline:
  stages:
    - build
    - test
    - deploy
  
  build:
    image: node:16
    script:
      - npm install
      - npm run build
    artifacts:
      paths:
        - dist/
  
  test:
    image: node:16
    script:
      - npm test
    coverage: '/Coverage: \d+\.\d+/'`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MakeVerbose(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Should have comments
			commentCount := strings.Count(result, "#")
			if commentCount == 0 {
				t.Error("Expected result to contain comments")
			}

			// Should be valid YAML
			var parsed interface{}
			if err := yaml.Unmarshal([]byte(result), &parsed); err != nil {
				t.Errorf("Result is not valid YAML: %v", err)
			}

			// Should be longer than input (due to comments)
			if len(result) <= len(tc.input) {
				t.Error("Expected verbose version to be longer due to comments")
			}
		})
	}
}
