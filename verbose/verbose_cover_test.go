package verbose

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestGetScalarTypeCommentCoverage(t *testing.T) {
	nodeLongString := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "this is a very long string that should definitely exceed fifty characters to trigger the string (long) branch",
	}
	optsDefault := Options{}
	if got := getScalarTypeComment(nodeLongString, optsDefault); got != "string (long)" {
		t.Errorf("getScalarTypeComment() = %v, want %v", got, "string (long)")
	}

	nodeEmptyString := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "",
	}
	optsExamples := Options{AddExamples: true}
	if got := getScalarTypeComment(nodeEmptyString, optsExamples); got != "string (empty)" {
		t.Errorf("getScalarTypeComment() = %v, want %v", got, "string (empty)")
	}
}

func TestProcessNodeDocument(t *testing.T) {
	node := &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "test"},
		},
	}
	opts := Options{Indent: 4, AddStructureComments: true}
	var buf bytes.Buffer
	if err := processNode(&buf, node, 0, "", opts); err != nil {
		t.Fatal(err)
	}
	if expected := "# YAML Document\ntest\n"; buf.String() != expected {
		t.Errorf("processNode() = %q; want %q", buf.String(), expected)
	}
}

func TestIsIntegerAndFloat(t *testing.T) {
	if isInteger("abc") {
		t.Errorf("Expected false for abc")
	}
	if isFloat("abc") {
		t.Errorf("Expected false for abc")
	}
}

func TestGetScalarTypeComment(t *testing.T) {
	nodeLong := &yaml.Node{Value: "this is a very long string that should be longer than 50 characters so it can trigger the long string path"}
	if getScalarTypeComment(nodeLong, DefaultOptions()) != "string (long)" {
		t.Errorf("Expected string (long)")
	}

	nodeEmpty := &yaml.Node{Value: ""}
	opts := DefaultOptions()
	opts.AddExamples = true
	if getScalarTypeComment(nodeEmpty, opts) != "string (empty)" {
		t.Errorf("Expected string (empty)")
	}
}

func TestProcessNodeNil(t *testing.T) {
	var b bytes.Buffer
	err := processNode(&b, nil, 0, "", DefaultOptions())
	if err != nil {
		t.Errorf("Expected no error for nil node")
	}
}

func TestVerboseErrorHandling(t *testing.T) {
	_, err := MakeVerboseWithOptions("invalid: \n  - yaml: [unclosed", DefaultOptions())
	if err == nil {
		t.Errorf("Expected error for invalid yaml")
	}

	err = validateIndentation("\tfoo: bar")
	if err == nil {
		t.Errorf("Expected error for tab indentation")
	}
}
