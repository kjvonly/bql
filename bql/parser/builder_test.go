package parser_test

import (
	"testing"

	"launchpad.net/kjvonly-bql/bql/parser"
	"launchpad.net/kjvonly-bql/bql/state"
)

func TestBuilderAddExpression(t *testing.T) {
	b := parser.NewBuilder(nil)
	e := b.AddExpression()

	if e.Expressions != nil {
		t.Fatalf("Should have nil Next")
	}

	if b.Expression == nil {
		t.Fatalf("Should have non nil Builder.Marker")
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

func TestBuilderAssignOrphanedExpressions(t *testing.T) {
	b := parser.NewBuilder(state.BQLLexer("="))
	b.AddExpression()

	b.OrphanedExpressions = append(b.OrphanedExpressions, []*parser.Expression{{}, {}}...)
	e := &parser.Expression{}
	b.AssignOrphanedExpressions(e)

	if len(b.OrphanedExpressions) != 0 {
		t.Fatalf("expected OrphanedExpressions to be 0")
	}
}
