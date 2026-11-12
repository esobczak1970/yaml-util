package verbose

import (
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestVerbosePreservesQuotedKeysAndValuesWithoutYAMLInjection(t *testing.T) {
	input := `"role: admin": "value\"\nadmin: true"
"line\nadmin: true": "second\nline"
single: 'it''s still safe'
path: "C:\\\\windows\\\\system32"
items:
  - "payload\"\n- unauthorized: true"
  - 'a''b'
`
	var want map[string]any
	if err := yaml.Unmarshal([]byte(input), &want); err != nil {
		t.Fatalf("invalid test fixture: %v", err)
	}

	output, err := MakeVerbose(input)
	if err != nil {
		t.Fatalf("MakeVerbose() failed: %v", err)
	}
	var got map[string]any
	if err := yaml.Unmarshal([]byte(output), &got); err != nil {
		t.Fatalf("transformed YAML must remain parseable: %v\n%s", err, output)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("verbose output changed YAML data\nwant: %#v\ngot: %#v\noutput:\n%s", want, got, output)
	}
	if _, injected := got["admin"]; injected {
		t.Fatalf("unexpected injected top-level key: %#v", got)
	}
	if _, injected := got["unauthorized"]; injected {
		t.Fatalf("unexpected injected top-level key: %#v", got)
	}
}

func TestVerbosePreservesAnchoredScalarsAndEscapedQuotes(t *testing.T) {
	input := `message: &msg "she said \"hello\""
copy: *msg
list:
  - &quote 'it''s safe'
  - *quote
`
	var want map[string]any
	if err := yaml.Unmarshal([]byte(input), &want); err != nil {
		t.Fatalf("invalid test fixture: %v", err)
	}
	output, err := MakeVerbose(input)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := yaml.Unmarshal([]byte(output), &got); err != nil {
		t.Fatalf("verbose output is invalid YAML: %v\n%s", err, output)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("anchor/alias semantics changed:\nwant: %#v\ngot: %#v\n%s", want, got, output)
	}
	if strings.Count(output, "&msg") != 1 || strings.Count(output, "&quote") != 1 {
		t.Fatalf("anchored scalars were emitted more than once:\n%s", output)
	}
}

func TestVerboseFormatMappingKeyEscapesYAMLMetacharacters(t *testing.T) {
	cases := []struct {
		value string
		style yaml.Style
		want  string
	}{
		{"plain", 0, "plain"},
		{"unsafe: key", 0, `"unsafe: key"`},
		{"a\"b", yaml.DoubleQuotedStyle, `"a\\"b"`},
		{"it's", yaml.SingleQuotedStyle, "'it''s'"},
		{"line\nnext: true", 0, `"line\\nnext: true"`},
	}
	for _, tc := range cases {
		got := formatMappingKey(&yaml.Node{Kind: yaml.ScalarNode, Style: tc.style, Value: tc.value})
		if got != tc.want {
			t.Errorf("formatMappingKey(%q) = %q, want %q", tc.value, got, tc.want)
		}
	}
}
