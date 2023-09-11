package parser_test

import (
	"testing"

	"launchpad.net/kjvonly-bql/bql/parser"
)

func TestBuilderMark(t *testing.T) {
	b := parser.Builder{}
	_ = b.Mark()

	if len(b.Marks) != 1 {
		t.Fatalf("Should have 1 mark but has %d marks", len(b.Marks))
	}
}

func TestParseFieldName(t *testing.T) {

	b := parser.Builder{}
	_ = b.Mark()

	if len(b.Marks) != 1 {
		t.Fatalf("Should have 1 mark but has %d marks", len(b.Marks))
	}
}

func TestAdvanceIfMatches(t *testing.T) {
	p := parser.Parser{}
	b := parser.Builder{}

	matches := p.AdvanceIfMatches(b, parser.VALID_FIELD_NAMES)

	if !matches {
		t.Fatalf("should match")
	}
}
