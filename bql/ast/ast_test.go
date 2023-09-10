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
	query := "book = jonn"
	_ = state.BQLLexer(query)

	//ast.ParseQuery(tokens)

}
