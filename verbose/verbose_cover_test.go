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

func TestProcessSequenceNodeErrorPropagation(t *testing.T) {
	node := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{
		{Kind: yaml.AliasNode, Value: "a"},
	}}
	var b bytes.Buffer
	// It doesn't actually error out, just checking it doesn't panic.
	_ = processSequenceNode(&b, node, 0, "", DefaultOptions())
}

func TestProcessMappingNodeErrorPropagation(t *testing.T) {
	node := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "k"},
		{Kind: yaml.AliasNode, Value: "a"},
	}}
	var b bytes.Buffer
	_ = processMappingNode(&b, node, 0, "", DefaultOptions())
}

func TestFormatScalarValueStyles(t *testing.T) {
	nSingle := &yaml.Node{Kind: yaml.ScalarNode, Style: yaml.SingleQuotedStyle, Value: "v"}
	if formatScalarValue(nSingle) != "'v'" {
		t.Errorf("Expected 'v'")
	}
}

// Missing test: processNode for DocumentNode that errors
func TestProcessNodeDocError(t *testing.T) {
	node := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{{Kind: yaml.AliasNode, Value: "a"}}}
	var b bytes.Buffer
	_ = processNode(&b, node, 0, "", DefaultOptions())
}

func TestMakeVerboseWithOptionsError(t *testing.T) {
	// A syntax error that triggers Unmarshal to fail
	_, err := MakeVerboseWithOptions("invalid:\n  - yaml: [unclosed", DefaultOptions())
	if err == nil {
		t.Errorf("Expected err")
	}
}

// A quick formatting fix for the go file.

// A test to cover processSequenceNode yaml.SequenceNode and yaml.MappingNode branches
func TestProcessSequenceNodeComplex(t *testing.T) {
	node := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{
		{Kind: yaml.SequenceNode, Content: []*yaml.Node{}},
		{Kind: yaml.MappingNode, Content: []*yaml.Node{}},
	}}
	var b bytes.Buffer
	_ = processSequenceNode(&b, node, 0, "", DefaultOptions())
}

// A test to cover MakeVerboseWithOptions error from processNode (although it doesnt error if valid YAML but the error branch exists in MakeVerboseWithOptions if processNode errors)
// Actually processNode only errors if mapping or sequence process errors, but they don't natively error unless they write to buffer and buffer fails (which bytes.Buffer doesn't)

func TestMakeVerboseWithOptions_OptionsLogic(t *testing.T) {
	opts := Options{AddTypeComments: false, AddStructureComments: false, AddExamples: false, Indent: 4}
	res, err := MakeVerboseWithOptions("k: v", opts)
	if err != nil {
		t.Fatal(err)
	}
	expected := "k: v\n"
	if res != expected {
		t.Errorf("Expected %q, got %q", expected, res)
	}

	optsTypeOnly := Options{AddTypeComments: true, AddStructureComments: false, AddExamples: false, Indent: 4}
	res2, err := MakeVerboseWithOptions("k: v", optsTypeOnly)
	if err != nil {
		t.Fatal(err)
	}
	expected2 := "k: v # string\n"
	if res2 != expected2 {
		t.Errorf("Expected %q, got %q", expected2, res2)
	}
}
