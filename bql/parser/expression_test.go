package parser_test

import (
	"testing"

	"launchpad.net/kjvonly-bql/bql/parser"
	"launchpad.net/kjvonly-bql/bql/state"
)

func TestDoneSuccess(t *testing.T) {
	b := parser.NewBuilder(state.BQLLexer("book = john"))
	e1 := b.AddExpression()
	e2 := b.AddExpression()

	e2.Done(state.EQ)
	e1.Done(state.IDENTIFIER)
}

func TestDoneFailure(t *testing.T) {
	b := parser.NewBuilder(state.BQLLexer("book = john"))
	e1 := b.AddExpression()

	e2 := b.AddExpression()
	e1.Expressions = append(e1.Expressions, e2)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	// second expression never called done
	e1.Done(state.IDENTIFIER)
}
