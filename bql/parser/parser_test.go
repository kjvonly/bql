package parser_test

import (
	"testing"

	"launchpad.net/kjvonly-bql/bql/parser"
	"launchpad.net/kjvonly-bql/bql/state"
)

// ///////////////////////////////////
// /////////// Builder ///////////////
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

	if token != state.STRING_LITERAL {
		t.Fatalf("expected %s but got %s", state.STRING_LITERAL, token)
	}
}

func TestBuilderAdvanceLexer(t *testing.T) {
	l := state.BQLLexer("book = john")
	b := parser.NewBuilder(l)

	b.AdvanceLexer()

	if b.CurrentToken.Type != state.STRING_LITERAL {
		t.Fatalf("expected current type to be STRING_LITERAL")
	}
}

// ////////////////////////////////////
// //////////// PARSER ///////////////
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

	input := "book = john"
	b.Lexer = state.BQLLexer(input)

	b.Mark()

	matches := p.AdvanceIfMatches(b, state.VALID_FIELD_NAMES)

	if !matches {
		t.Fatalf("should match")
	}
}
