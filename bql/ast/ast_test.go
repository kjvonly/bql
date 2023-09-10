package ast_test

import (
	"testing"

	"launchpad.net/kjvonly-bql/bql/ast"
)

func TestSimpleExpression(t *testing.T) {
	query := "book = jonn"

	es := ast.EqualStmt{
		Field: ast.Ident{Name: "book"},
		Expr:  &ast.Ident{Name: "John"},
	}

	tree := ast.Ast{}

	tree.Walk()

}
