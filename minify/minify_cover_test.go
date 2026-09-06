package minify

import (
	"testing"
	"strings"
	"gopkg.in/yaml.v3"
)

func TestEmitNodeDocument(t *testing.T) {
	var b strings.Builder
	node := &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "test"},
		},
	}
	if err := emitNode(&b, node, 0, true); err != nil {
		t.Fatal(err)
	}
	if b.String() != "test" {
		t.Fatalf("expected 'test', got '%s'", b.String())
	}
}
