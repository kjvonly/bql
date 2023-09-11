package ast_test

import (
	"testing"

	"launchpad.net/kjvonly-bql/bql/ast"
	"launchpad.net/kjvonly-bql/bql/state"
)

func TestInvalidEqualStmt(t *testing.T) {
	query := "="
	tokens := state.BQLLexer(query)
	a := ast.Ast{}
	err := a.Generate(tokens)

	if err == nil {
		t.Fatalf("expected query error")
	}
}

func TestValid(t *testing.T) {
	query := "book = john"

	a := ast.Ast{}
	a.Tokens = state.BQLLexer(query)
	err := a.ParseQuery()

	if err != nil {
		t.Fatalf("did not expect query error: error %s", err)
	}

	if len(a.Elements) != 3 {
		t.Fatalf("expected 3 element but had %d", len(a.Elements))
	}
}

func TestFieldNotFirstElement(t *testing.T) {
	query := "john = love"

	a := ast.Ast{}
	a.Tokens = state.BQLLexer(query)
	err := a.ParseQuery()

	if err == nil {
		t.Fatalf("did expect query error: %s", err)
	}
}

func TestAnd(t *testing.T) {
	query := "book = john and"

	a := ast.Ast{}
	a.Tokens = state.BQLLexer(query)
	err := a.ParseQuery()

	if err != nil {
		t.Fatalf("did not expect query error: error %s", err)
	}

	if len(a.Elements) != 4 {
		t.Fatalf("expected 4 element but had %d", len(a.Elements))
	}
}
