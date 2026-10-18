package maxify

import (
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
