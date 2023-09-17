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

	m2 := b.Mark()
	m1.Children = append(m1.Children, m2)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	// second mark never called done
	m1.Done(state.IDENTIFIER)
}

func TestMarkerPrecede(t *testing.T) {
	b := parser.NewBuilder(state.BQLLexer("book = john"))
	m1 := b.Mark()
	m1.Precede(b)

	if b.Markers.Tail != m1 {
		t.Fatalf("expected new tail to be preceded marker but was %+v", b.Markers.Tail)
	}
}

// ///////////////////////////////////
// /////////// Builder ///////////////
func TestBuilderMark(t *testing.T) {
	b := parser.NewBuilder(nil)
	m := b.Mark()

	if m.Children != nil && m.Parent != nil {
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
	b.Mark()

	b.OrphanedChildren = append(b.OrphanedChildren, []*parser.Marker{{}, {}}...)
	m := &parser.Marker{}
	b.AssignOrphanedChildren(m)

	for _, c := range m.Children {
		if c.Parent != m {
			t.Fatalf("expected children to have correct parent")
		}
	}

	if len(b.OrphanedChildren) != 0 {
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

	if b.Markers.Tail.Type != state.IDENTIFIER {
		t.Fatalf("expected marker have IDENTIFIER type")
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
	queryMarker := b.Mark() // Query Marker needed for testing
	success := p.ParseOrClause(b)
	queryMarker.Done(state.QUERY)
	b.AssignOrphanedChildren(queryMarker)

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

	ma := flattenMarkers(b.Markers.Head)
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
	queryMarker := b.Mark() // Query Marker needed for testing
	success := p.ParseOrClause(b)
	queryMarker.Done(state.QUERY)
	b.AssignOrphanedChildren(queryMarker)

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

	ma := flattenMarkers(b.Markers.Head)
	for i := 0; i < len(ma); i++ {
		if expectedElementTypeOrdered[i] != ma[i].Type {
			t.Fatalf("expected type %s but got %s", expectedElementTypeOrdered[i], ma[i].Type)
		}
	}
}

func TestParseAndClauseShouldNotSucceed(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book"))

	queryMarker := b.Mark()
	success := p.ParseAndClause(b)
	queryMarker.Done(state.QUERY)

	if success {
		t.Fatalf("expected not to succeed")
	}

}

func TestParseAndClauseShouldSucceed(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book = john and book = mark and book = matthew"))
	b.AdvanceLexer()
	queryMarker := b.Mark() // Query Marker needed for testing
	success := p.ParseAndClause(b)
	queryMarker.Done(state.QUERY)

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

	ma := flattenMarkers(b.Markers.Head)
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

func flattenMarkers(m *parser.Marker) []*parser.Marker {

	ma := []*parser.Marker{}

	ma = append(ma, m)

	for i := 0; i < len(m.Children); i++ {
		ma = append(ma, flattenMarkers(m.Children[i])...)
	}
	return ma
}

func TestParseTerminalClauseProperFieldName(t *testing.T) {
	p := parser.Parser{}
	b := parser.NewBuilder(state.BQLLexer("book = john"))
	queryMarker := b.Mark() // Query marker
	b.AdvanceLexer()
	success := p.ParseTerminalClause(b)
	b.AssignOrphanedChildren(queryMarker)
	queryMarker.Done(state.QUERY)

	if !success {
		t.Fatalf("expected true")
	}

	expectedElementTypeOrdered := []state.ElementType{state.QUERY, state.SIMPLE_CLAUSE, state.IDENTIFIER, state.LITERAL}

	ma := flattenMarkers(b.Markers.Head)
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

	if b.Markers.Tail.Type != state.LITERAL {
		t.Fatalf("expected Literal marker but got %s", b.Markers.Tail.Type)
	}

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
