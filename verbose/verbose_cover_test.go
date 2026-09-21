package verbose

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"testing"
)

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
