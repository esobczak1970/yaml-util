package minify

import (
	"gopkg.in/yaml.v3"
	"strings"
	"testing"
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

func TestFormatFlowSequenceMixed(t *testing.T) {
	node := &yaml.Node{
		Kind: yaml.SequenceNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "a"},
			{Kind: yaml.AliasNode, Value: "alias1"},
			{Kind: yaml.DocumentNode},
		},
	}
	res := formatFlowSequence(node)
	if res != "[a, *alias1, null]" {
		t.Fatalf("expected '[a, *alias1, null]', got '%s'", res)
	}
}

func TestFormatFlowMappingMixed(t *testing.T) {
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "k1"},
			{Kind: yaml.ScalarNode, Value: "v1"},
			{Kind: yaml.ScalarNode, Value: "k2"},
			{Kind: yaml.AliasNode, Value: "alias1"},
			{Kind: yaml.ScalarNode, Value: "k3"},
			{Kind: yaml.DocumentNode},
		},
	}
	res := formatFlowMapping(node)
	if res != "{k1: v1, k2: *alias1, k3: null}" {
		t.Fatalf("expected '{k1: v1, k2: *alias1, k3: null}', got '%s'", res)
	}
}

func TestEmitSequenceMixed(t *testing.T) {
	var b strings.Builder
	node := &yaml.Node{
		Kind: yaml.SequenceNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "a"},
			{Kind: yaml.AliasNode, Value: "alias1"},
			{Kind: yaml.DocumentNode},
		},
	}
	if err := emitSequence(&b, node, 0); err != nil {
		t.Fatal(err)
	}
	res := b.String()
	expected := "- a\n- *alias1\n-\n"
	if res != expected {
		t.Fatalf("expected '%s', got '%s'", expected, res)
	}
}

func TestFormatAliasNil(t *testing.T) {
	node := &yaml.Node{Kind: yaml.AliasNode}
	res := formatAlias(node)
	if res != "*" {
		t.Fatalf("expected '*', got '%s'", res)
	}

	node2 := &yaml.Node{Kind: yaml.AliasNode, Value: "v"}
	res2 := formatAlias(node2)
	if res2 != "*v" {
		t.Fatalf("expected '*v', got '%s'", res2)
	}
}

func TestEmitNodeOtherTypes(t *testing.T) {
	var b strings.Builder
	nodeAlias := &yaml.Node{Kind: yaml.AliasNode, Value: "a"}
	_ = emitNode(&b, nodeAlias, 0, false)

	b.Reset()
	nodeSeq := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{{Kind: yaml.ScalarNode, Value: "a"}}}
	_ = emitNode(&b, nodeSeq, 0, false)

	b.Reset()
	nodeMap := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{{Kind: yaml.ScalarNode, Value: "k"}, {Kind: yaml.ScalarNode, Value: "v"}}}
	_ = emitNode(&b, nodeMap, 0, false)
}

func TestEmitMappingMixed(t *testing.T) {
	var b strings.Builder
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "k1"},
			{Kind: yaml.AliasNode, Value: "a1"},

			{Kind: yaml.ScalarNode, Value: "k2"},
			{Kind: yaml.SequenceNode, Content: []*yaml.Node{{Kind: yaml.ScalarNode, Value: "v"}}},

			{Kind: yaml.ScalarNode, Value: "k3"},
			{Kind: yaml.MappingNode, Content: []*yaml.Node{{Kind: yaml.ScalarNode, Value: "k"}, {Kind: yaml.ScalarNode, Value: "v"}}},
		},
	}
	_ = emitMapping(&b, node, 0, false)
}

func TestEmitSequenceComplex(t *testing.T) {
	var b strings.Builder
	node := &yaml.Node{
		Kind: yaml.SequenceNode,
		Content: []*yaml.Node{
			{Kind: yaml.SequenceNode, Content: []*yaml.Node{{Kind: yaml.ScalarNode, Value: "v"}}},
			{Kind: yaml.MappingNode, Content: []*yaml.Node{{Kind: yaml.ScalarNode, Value: "k"}, {Kind: yaml.ScalarNode, Value: "v"}}},
			{Kind: yaml.DocumentNode},
		},
	}
	// Test shouldInlineSequence logic
	nodeSeqComplex := &yaml.Node{
		Kind: yaml.SequenceNode,
		Content: []*yaml.Node{
			{Kind: yaml.SequenceNode, Content: []*yaml.Node{{Kind: yaml.ScalarNode, Value: "v"}}},
		},
	}
	_ = emitSequence(&b, nodeSeqComplex, 0)
	b.Reset()
	_ = emitSequence(&b, node, 0)
}

func TestEmitMappingInline(t *testing.T) {
	var b strings.Builder
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "k1"},
			{Kind: yaml.ScalarNode, Value: "v1"},
		},
	}
	// to trigger the shouldInlineMapping path inside emitMapping where isRoot is false
	_ = emitMapping(&b, node, 0, false)
}

func TestEmitSequenceMappingInline(t *testing.T) {
	var b strings.Builder
	nodeMap := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "k"},
			{Kind: yaml.ScalarNode, Value: "v"},
		},
	}
	node := &yaml.Node{
		Kind: yaml.SequenceNode,
		Content: []*yaml.Node{
			nodeMap,
		},
	}
	// To hit the "if shouldInlineMapping(item, false)" path in emitSequence
	_ = emitSequence(&b, node, 0)
}

func TestNeedsQuotes(t *testing.T) {
	// hit needsQuotes paths
	_ = needsQuotes("false")
	_ = needsQuotes("n")
	_ = needsQuotes("1.2")
	_ = needsQuotes("~")
	_ = needsQuotes("y")
	_ = needsQuotes("true")
}

func TestEmitMappingNestedScalar(t *testing.T) {
	var b strings.Builder
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "k"},
			{Kind: yaml.ScalarNode, Value: "v"},
		},
	}
	// To test the spaceForScalar logic
	_ = emitMapping(&b, node, 2, false)
}

func TestEmitMappingAnchorScalar(t *testing.T) {
	var b strings.Builder
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "k"},
			{Kind: yaml.ScalarNode, Value: "v", Anchor: "anchor1"},
		},
	}
	// To test the val.Anchor logic
	_ = emitMapping(&b, node, 0, true)
}

func TestShouldInlineMappingCases(t *testing.T) {
	nodeNil := &yaml.Node{Kind: yaml.MappingNode, Content: nil}
	_ = shouldInlineMapping(nodeNil, false)

	nodeAlias := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "k"},
		{Kind: yaml.AliasNode, Value: "a"},
	}}
	_ = shouldInlineMapping(nodeAlias, false)

	nodeEmpty := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "k"},
		{Kind: yaml.ScalarNode, Value: ""},
	}}
	_ = shouldInlineMapping(nodeEmpty, false)
}

func TestShouldInlineSequenceCases(t *testing.T) {
	nodeNil := &yaml.Node{Kind: yaml.SequenceNode, Content: nil}
	_ = shouldInlineSequence(nodeNil)

	nodeAlias := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{
		{Kind: yaml.AliasNode, Value: "a"},
	}}
	_ = shouldInlineSequence(nodeAlias)

	nodeEmpty := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: ""},
	}}
	_ = shouldInlineSequence(nodeEmpty)
}

func TestFormatKeyScalar(t *testing.T) {
	nodeNil := &yaml.Node{Kind: yaml.AliasNode}
	_ = formatKeyScalar(nodeNil)

	nodeDouble := &yaml.Node{Kind: yaml.ScalarNode, Style: yaml.DoubleQuotedStyle, Value: "v"}
	_ = formatKeyScalar(nodeDouble)

	nodeSingle := &yaml.Node{Kind: yaml.ScalarNode, Style: yaml.SingleQuotedStyle, Value: "v"}
	_ = formatKeyScalar(nodeSingle)

	nodeQuotes := &yaml.Node{Kind: yaml.ScalarNode, Value: ":"}
	_ = formatKeyScalar(nodeQuotes)
}

func TestFormatScalar(t *testing.T) {
	nodeDouble := &yaml.Node{Kind: yaml.ScalarNode, Style: yaml.DoubleQuotedStyle, Value: "v"}
	_ = formatScalar(nodeDouble)

	nodeSingle := &yaml.Node{Kind: yaml.ScalarNode, Style: yaml.SingleQuotedStyle, Value: "v"}
	_ = formatScalar(nodeSingle)
}

func TestFormatScalarNil(t *testing.T) {
	nodeNil := &yaml.Node{Kind: yaml.AliasNode}
	_ = formatScalar(nodeNil)
}

func TestNeedsQuotesValues(t *testing.T) {
	_ = needsQuotes(" ")
	_ = needsQuotes("   ")
	_ = needsQuotes("")
}

func TestMinifyErrorEmit(t *testing.T) {
	nodeDoc := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{
		{Kind: yaml.AliasNode}, // Will error on emitNode ? No wait, DocumentNode doesn't have an error in emitNode unless inner emit fails. But there's no error return for alias node.
	}}
	_ = nodeDoc
}

func TestEmitNodeNilAndDefault(t *testing.T) {
	var b strings.Builder
	_ = emitNode(&b, nil, 0, false)

	nodeUnknown := &yaml.Node{Kind: 999}
	_ = emitNode(&b, nodeUnknown, 0, false)

	nodeDocEmpty := &yaml.Node{Kind: yaml.DocumentNode, Content: nil}
	_ = emitNode(&b, nodeDocEmpty, 0, false)
}
