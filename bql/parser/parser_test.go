package parser_test

import (
	"testing"

	"launchpad.net/kjvonly-bql/bql/parser"
	"launchpad.net/kjvonly-bql/bql/state"
)

func TestDoneSuccess(t *testing.T) {
	b := parser.NewBuilder(state.BQLLexer("book = john"))
	m1 := b.Mark()
	m2 := b.Mark()

	m2.Done(state.EQ)
	m1.Done(state.IDENTIFIER)
}

func TestDoneFailure(t *testing.T) {
	b := parser.NewBuilder(state.BQLLexer("book = john"))
	m1 := b.Mark()
	b.Mark()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	// second mark never called done
	m1.Done(state.IDENTIFIER)
}

// ///////////////////////////////////
// /////////// Builder ///////////////
func TestBuilderMark(t *testing.T) {
	b := parser.NewBuilder(nil)
	m := b.Mark()

	if m.Next != nil && m.Prev != nil {
		t.Fatalf("Should have nil Next and Prev")
	}

	if b.Markers == nil {
		t.Fatalf("Should have non nil Builder.Marker")
	}

	if b.Markers.Head != b.Markers.Tail {
		t.Fatalf("Should assign head to tail")
	}

	if m != b.Markers.Head {
		t.Fatalf("Should return head marker")
	}

	m2 := b.Mark()

	if m2 != b.Markers.Tail {
		t.Fatalf("Should append marker to end of linked list")
	}

	if b.Markers.Head.Next != m2 {
		t.Fatalf("Should assign next to previous tail marker ")
	}

	if m2.Prev != b.Markers.Head {
		t.Fatalf("Should assign Prev to previous marker")
	}
}

func TestBuilderGetTokenType(t *testing.T) {
	l := state.BQLLexer("book = john")
	b := parser.NewBuilder(l)
	b.AdvanceLexer()

	token := b.GetTokenType()

	if token != state.IDENTIFIER {
		t.Fatalf("expected %s but got %s", state.IDENTIFIER, token)
	}
}

func TestBuilderGetTokenTypeEmptyString(t *testing.T) {
	l := state.BQLLexer("book")
	b := parser.NewBuilder(l)

	token := b.GetTokenType()

	if token != "" {
		t.Fatalf("expected empty string but got %s", token)
	}
}

func TestBuilderAdvanceLexer(t *testing.T) {
	l := state.BQLLexer("book = john")
	b := parser.NewBuilder(l)

	b.AdvanceLexer()

	if b.CurrentToken.Type != state.IDENTIFIER {
		t.Fatalf("expected current type to be %s", state.IDENTIFIER)
	}
}

func TestBuilderError(t *testing.T) {
	b := &parser.Builder{}

	b.Error("expecting field name")

	// TODO finish test
}

// ////////////////////////////////////
// //////////// PARSER ///////////////
func TestParseInvalidFieldName(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("="))
	b.AdvanceLexer()
	success := p.ParseFieldName(b)
	if !b.Markers.Head.IsDropped {
		t.Fatalf("expected mark to have been dropped")
	}
	if success {
		t.Fatalf("expected parseFieldName to have failed")
	}
}

func TestParseValidFieldName(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book"))
	b.AdvanceLexer()
	success := p.ParseFieldName(b)
	if b.Markers.Head.IsDropped {
		t.Fatalf("expected mark not to have been dropped")
	}
	if !success {
		t.Fatalf("expected parseFieldName to have succeeded")
	}

	if b.Markers.Tail.Type != state.IDENTIFIER {
		t.Fatalf("expected marker have IDENTIFIER type")
	}
}

func TestParseTerminalClauseNotProperFieldName(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("="))
	b.AdvanceLexer()
	success := p.ParseTerminalClause(b)

	if success {
		t.Fatalf("expected false")
	}
}

func TestParseTerminalClauseProperFieldName(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book"))
	b.AdvanceLexer()
	success := p.ParseTerminalClause(b)

	if !success {
		t.Fatalf("expected false")
	}
}
func TestAdvanceIfMatches(t *testing.T) {
	p := parser.Parser{}

	input := "book"
	b := parser.NewBuilder(state.BQLLexer(input))

	b.AdvanceLexer()

	matches := p.AdvanceIfMatches(b, state.VALID_FIELD_NAMES)

	if !matches {
		t.Fatalf("should match")
	}
}
