package parser_test

import (
	"testing"

	"launchpad.net/kjvonly-bql/bql/parser"
	"launchpad.net/kjvonly-bql/bql/state"
)

func TestParseInvalidFieldName(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("="))
	b.AdvanceLexer()
	success := p.ParseFieldName(b)
	if success {
		t.Fatalf("expected parseFieldName to have failed")
	}
}

func TestParseValidFieldName(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book"))
	b.AdvanceLexer()
	success := p.ParseFieldName(b)

	if !success {
		t.Fatalf("expected parseFieldName to have succeeded")
	}
}

func TestParseOrClauseShouldNotSucceed(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book"))

	success := p.ParseOrClause(b)

	if success {
		t.Fatalf("expected not to succeed")
	}
}

func flattenExpressions(m *parser.Expression) []*parser.Expression {

	ma := []*parser.Expression{}

	ma = append(ma, m)

	for i := 0; i < len(m.Expressions); i++ {
		ma = append(ma, flattenExpressions(m.Expressions[i])...)
	}
	return ma
}

func TestParseQueryShouldSucceed(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book = john and book = mark or book = matthew"))
	b.AdvanceLexer()
	success := p.ParseQuery(b)

	if !success {
		t.Fatalf("expected to succeed")
	}

	expectedExpressionTypeOrdered := []state.ElementType{
		state.QUERY,
		state.OR_CLAUSE,
		state.AND_CLAUSE,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
	}

	es := flattenExpressions(b.Expression)
	for i := 0; i < len(es); i++ {
		if expectedExpressionTypeOrdered[i] != es[i].Type {
			t.Fatalf("expected type %s but got %s", expectedExpressionTypeOrdered[i], es[i].Type)
		}
	}
}

func TestParseOrClauseShouldSucceed(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book = john or book = mark or book = matthew"))
	b.AdvanceLexer()
	success := p.ParseOrClause(b)
	b.Expression.Done(state.QUERY)
	b.AssignOrphanedExpressions(b.Expression)

	if !success {
		t.Fatalf("expected to succeed")
	}

	expectedExpressionTypeOrdered := []state.ElementType{
		state.QUERY,
		state.OR_CLAUSE,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
	}

	es := flattenExpressions(b.Expression)
	for i := 0; i < len(es); i++ {
		if expectedExpressionTypeOrdered[i] != es[i].Type {
			t.Fatalf("expected type %s but got %s", expectedExpressionTypeOrdered[i], es[i].Type)
		}
	}
}

func TestParseAndOrClauseShouldSucceed(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book = john or book = mark and book = matthew"))
	b.AdvanceLexer()
	success := p.ParseOrClause(b)
	b.Expression.Done(state.QUERY)
	b.AssignOrphanedExpressions(b.Expression)

	if !success {
		t.Fatalf("expected to succeed")
	}

	expectedExpressionTypeOrdered := []state.ElementType{
		state.QUERY,
		state.OR_CLAUSE,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
		state.AND_CLAUSE,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
	}

	es := flattenExpressions(b.Expression)
	for i := 0; i < len(es); i++ {
		if expectedExpressionTypeOrdered[i] != es[i].Type {
			t.Fatalf("expected type %s but got %s", expectedExpressionTypeOrdered[i], es[i].Type)
		}
	}
}

func TestParseAndClauseShouldNotSucceed(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book"))

	success := p.ParseAndClause(b)
	b.Expression.Done(state.QUERY)
	b.AssignOrphanedExpressions(b.Expression)

	if success {
		t.Fatalf("expected not to succeed")
	}

}

func TestParseAndClauseShouldSucceed(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book = john and book = mark and book = matthew"))
	b.AdvanceLexer()
	success := p.ParseAndClause(b)
	b.Expression.Done(state.QUERY)
	b.AssignOrphanedExpressions(b.Expression)

	if !success {
		t.Fatalf("expected to succeed")
	}

	expectedExpressionTypeOrdered := []state.ElementType{
		state.QUERY,
		state.AND_CLAUSE,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
		state.SIMPLE_CLAUSE,
		state.IDENTIFIER,
		state.LITERAL,
	}

	es := flattenExpressions(b.Expression)
	for i := 0; i < len(es); i++ {
		if expectedExpressionTypeOrdered[i] != es[i].Type {
			t.Fatalf("expected type %s but got %s", expectedExpressionTypeOrdered[i], es[i].Type)
		}
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
	b := parser.NewBuilder(state.BQLLexer("book = john"))
	b.AdvanceLexer()
	success := p.ParseTerminalClause(b)
	b.Expression.Done(state.QUERY)
	b.AssignOrphanedExpressions(b.Expression)

	if !success {
		t.Fatalf("expected true")
	}

	expectedExpressionTypeOrdered := []state.ElementType{state.QUERY, state.SIMPLE_CLAUSE, state.IDENTIFIER, state.LITERAL}

	es := flattenExpressions(b.Expression)
	for i := 0; i < len(expectedExpressionTypeOrdered); i++ {
		if expectedExpressionTypeOrdered[i] != es[i].Type {
			t.Fatalf("expected type %s but got %s", expectedExpressionTypeOrdered[i], es[i].Type)
		}
	}
}

func TestParseOperand(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("john"))
	b.AdvanceLexer()
	parsed := p.ParseOperand(b)

	if !parsed {
		t.Fatalf("expected parsed to be true but was false")
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
