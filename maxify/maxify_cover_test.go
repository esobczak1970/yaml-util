package maxify

import (
	"gopkg.in/yaml.v3"
	"testing"
)

func TestEnsureTrailingNewlineEmpty(t *testing.T) {
	if got := ensureTrailingNewline(""); got != "" {
		t.Errorf("ensureTrailingNewline(\"\") = %q; want %q", got, "")
	}
}

func TestEnsureTrailingNewline(t *testing.T) {
	if got := ensureTrailingNewline("a"); got != "a\n" {
		t.Errorf("ensureTrailingNewline(\"a\") = %q; want %q", got, "a\n")
	}
	if got := ensureTrailingNewline("a\n"); got != "a\n" {
		t.Errorf("ensureTrailingNewline(\"a\\n\") = %q; want %q", got, "a\n")
	}
}

func TestMaxifyWithOptionsEmpty(t *testing.T) {
	res, err := MaxifyWithOptions("   \n", DefaultOptions())
	if err != nil || res != "" {
		t.Errorf("Expected empty result, got %q, %v", res, err)
	}
}

func TestMaxifyWithOptionsBasic(t *testing.T) {
	opts := DefaultOptions()
	opts.Indent = 2
	res, err := MaxifyWithOptions("key: value", opts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "key: value\n"
	if res != expected {
		t.Errorf("Expected %q, got %q", expected, res)
	}
}

func TestInsertSpaceAfterColon_NoSpace(t *testing.T) {
	in := "key:value"
	out := insertSpaceAfterColon(in)
	if out != "key: value" {
		t.Errorf("Expected 'key: value', got '%s'", out)
	}
}

func TestEnsureTrailingNewline_NoNewline(t *testing.T) {
	in := "foo: bar"
	out := ensureTrailingNewline(in)
	if out != "foo: bar\n" {
		t.Errorf("Expected 'foo: bar\\n', got '%s'", out)
	}
}

func TestMaxifyWithOptions_FormatError(t *testing.T) {
	_, err := MaxifyWithOptions("\tinvalid: indent", DefaultOptions())
	if err == nil {
		t.Errorf("Expected error for invalid indentation")
	}
}

func TestExpandScalar_Shorthand(t *testing.T) {
	node := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: 0,
		Value: "y",
	}
	expandScalar(node)
	if node.Value != "true" || node.Tag != "!!bool" {
		t.Errorf("Expected true and !!bool, got %s and %s", node.Value, node.Tag)
	}

	node = &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: 0,
		Value: "no",
	}
	expandScalar(node)
	if node.Value != "false" || node.Tag != "!!bool" {
		t.Errorf("Expected false and !!bool, got %s and %s", node.Value, node.Tag)
	}

	node = &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: 0,
		Value: "~",
	}
	expandScalar(node)
	if node.Value != "null" || node.Tag != "!!null" {
		t.Errorf("Expected null and !!null, got %s and %s", node.Value, node.Tag)
	}
}

func TestValidateIndentation_Tabs(t *testing.T) {
	err := validateIndentation("\tfoo: bar")
	if err == nil || err.Error() != "invalid indentation: tabs are not allowed" {
		t.Errorf("Expected error about tabs, got %v", err)
	}
}

func TestExpandNode_Nil(t *testing.T) {
	expandNode(nil) // should not panic
}
