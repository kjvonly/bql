package parser_test

import (
	"testing"

	"launchpad.net/kjvonly-bql/bql/parser"
	"launchpad.net/kjvonly-bql/bql/state"
)

func TestDoneSuccess(t *testing.T) {
	b := parser.NewBuilder(state.BQLLexer("book = john"))
	m1 := b.AddExpression()
	m2 := b.AddExpression()

	m2.Done(state.EQ)
	m1.Done(state.IDENTIFIER)
}

func TestDoneFailure(t *testing.T) {
	b := parser.NewBuilder(state.BQLLexer("book = john"))
	m1 := b.AddExpression()

	m2 := b.AddExpression()
	m1.Expressions = append(m1.Expressions, m2)

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
	m := b.AddExpression()

	if m.Expressions != nil && m.Parent != nil {
		t.Fatalf("Should have nil Next and Prev")
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

func TestBuilderAssignOrphanedChildren(t *testing.T) {
	b := parser.NewBuilder(state.BQLLexer("="))
	b.AddExpression()

	b.OrphanedExpressions = append(b.OrphanedExpressions, []*parser.Expression{{}, {}}...)
	m := &parser.Expression{}
	b.AssignOrphanedChildren(m)

	for _, c := range m.Expressions {
		if c.Parent != m {
			t.Fatalf("expected children to have correct parent")
		}
	}

	if len(b.OrphanedExpressions) != 0 {
		t.Fatalf("expected OrphanedChildren to be 0")
	}
}

// ////////////////////////////////////
// //////////// PARSER ///////////////
func TestParseInvalidFieldName(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("="))
	b.AdvanceLexer()
	success := p.ParseFieldName(b)
	// if !b.Markers.Head.IsDropped {
	// 	t.Fatalf("expected mark to have been dropped")
	// }
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

func TestParseOrClauseShouldSucceed(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book = john or book = mark or book = matthew"))
	b.AdvanceLexer()
	success := p.ParseOrClause(b)
	b.Expression.Done(state.QUERY)
	b.AssignOrphanedChildren(b.Expression)

	if !success {
		t.Fatalf("expected to succeed")
	}

	expectedElementTypeOrdered := []state.ElementType{
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

	ma := flattenMarkers(b.Expression)
	for i := 0; i < len(ma); i++ {
		if expectedElementTypeOrdered[i] != ma[i].Type {
			t.Fatalf("expected type %s but got %s", expectedElementTypeOrdered[i], ma[i].Type)
		}
	}
}

func TestParseAndOrClauseShouldSucceed(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book = john or book = mark and book = matthew"))
	b.AdvanceLexer()
	success := p.ParseOrClause(b)
	b.Expression.Done(state.QUERY)
	b.AssignOrphanedChildren(b.Expression)

	if !success {
		t.Fatalf("expected to succeed")
	}

	expectedElementTypeOrdered := []state.ElementType{
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

	ma := flattenMarkers(b.Expression)
	for i := 0; i < len(ma); i++ {
		if expectedElementTypeOrdered[i] != ma[i].Type {
			t.Fatalf("expected type %s but got %s", expectedElementTypeOrdered[i], ma[i].Type)
		}
	}
}

func TestParseAndClauseShouldNotSucceed(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book"))

	success := p.ParseAndClause(b)
	b.Expression.Done(state.QUERY)
	b.AssignOrphanedChildren(b.Expression)

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
	b.AssignOrphanedChildren(b.Expression)

	if !success {
		t.Fatalf("expected to succeed")
	}

	expectedElementTypeOrdered := []state.ElementType{
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

	ma := flattenMarkers(b.Expression)
	for i := 0; i < len(ma); i++ {
		if expectedElementTypeOrdered[i] != ma[i].Type {
			t.Fatalf("expected type %s but got %s", expectedElementTypeOrdered[i], ma[i].Type)
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

func flattenMarkers(m *parser.Expression) []*parser.Expression {

	ma := []*parser.Expression{}

	ma = append(ma, m)

	for i := 0; i < len(m.Expressions); i++ {
		ma = append(ma, flattenMarkers(m.Expressions[i])...)
	}
	return ma
}

func TestParseTerminalClauseProperFieldName(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book = john"))
	b.AdvanceLexer()
	success := p.ParseTerminalClause(b)
	b.Expression.Done(state.QUERY)
	b.AssignOrphanedChildren(b.Expression)

	if !success {
		t.Fatalf("expected true")
	}

	expectedElementTypeOrdered := []state.ElementType{state.QUERY, state.SIMPLE_CLAUSE, state.IDENTIFIER, state.LITERAL}

	ma := flattenMarkers(b.Expression)
	for i := 0; i < len(expectedElementTypeOrdered); i++ {
		if expectedElementTypeOrdered[i] != ma[i].Type {
			t.Fatalf("expected type %s but got %s", expectedElementTypeOrdered[i], ma[i].Type)
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
