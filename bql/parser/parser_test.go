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

func TestBuilderGetToken(t *testing.T) {
	b := parser.Builder{}
	_ = b.Mark()

	if len(b.Marks) != 1 {
		t.Fatalf("Should have 1 mark but has %d marks", len(b.Marks))
	}

	token := b.GetToken()

	if token != parser.STRING_LITERAL {
		t.Fatalf("expected %s but got %s", parser.STRING_LITERAL, token)
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

	b.Mark()

	matches := p.AdvanceIfMatches(b, parser.VALID_FIELD_NAMES)

	if !matches {
		t.Fatalf("should match")
	}
}
