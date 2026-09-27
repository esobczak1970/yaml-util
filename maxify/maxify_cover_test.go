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
