package maxify

import (
	"testing"
)

func TestEnsureTrailingNewlineEmpty(t *testing.T) {
	if got := ensureTrailingNewline(""); got != "" {
		t.Errorf("ensureTrailingNewline(\"\") = %q; want %q", got, "")
	}
}
