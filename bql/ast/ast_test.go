package ast_test

import (
	"testing"

	"launchpad.net/kjvonly-bql/bql/ast"
	"launchpad.net/kjvonly-bql/bql/state"
)

func TestSimpleExpression(t *testing.T) {
	query := "book = jonn"
	tokens := state.BQLLexer(query)

	// es := ast.EqualStmt{
	// 	Field: ast.Ident{Name: "book"},
	// 	Expr:  &ast.Ident{Name: "John"},
	// }

	tree := ast.Ast{}

	expectedExpr := tree.Generate(tokens)

	expectedExpr.Print()
	t.Fail()
}

func TestParseQuery(t *testing.T) {
	query := "="
	tokens := state.BQLLexer(query)
	a := ast.Ast{}
	err := a.Generate(tokens)

	if err == nil {
		t.Fatalf("expected query error")
	}

}

func TestInvalidEqualStmt(t *testing.T) {
	query := "book ="

	a := ast.Ast{}
	a.Tokens = state.BQLLexer(query)
	err := a.ParseQuery()

	if err != nil {
		t.Fatalf("did not expect query error: error %s", err)
	}

	if len(a.Elements) != 1 {
		t.Fatalf("expected 1 element but had %d", len(a.Elements))
	}

}
