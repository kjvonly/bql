package parser

import (
	"launchpad.net/kjvonly-bql/bql/state"
	"launchpad.net/kjvonly-bql/lex"
)

type Token struct {
	Type  state.ElementType
	Value interface{}
}

// NewExpression creates root Marker
func NewExpression() *Expression {
	return &Expression{}
}

func NewMarkerList() *MarkerList {
	return &MarkerList{}
}

type MarkerList struct {
	Head *Expression
	Tail *Expression
}

type Expression struct {
	Children []*Expression
	Parent   *Expression
	IsDone   bool
	Type     state.ElementType
	Value    interface{}
}

func checkAllExpressionsDone(n []*Expression) {
	for i := 0; i < len(n); i++ {
		if !n[i].IsDone {
			//TODO should change panic to something else
			panic("all markers past this marker not done.")
		}

		checkAllExpressionsDone(n[i].Children)
	}
}

func (m *Expression) Done(t state.ElementType) {
	checkAllExpressionsDone(m.Children)
	m.IsDone = true
	m.Type = t
}

type Builder struct {
	Markers          *Expression
	Lexer            *lex.Lexer
	CurrentToken     Token
	OrphanedChildren []*Expression
}

func NewBuilder(lex *lex.Lexer) *Builder {
	return &Builder{
		Lexer:   lex,
		Markers: &Expression{},
	}
}

// Mark adds a placeholder for new
func (b *Builder) Mark() *Expression {
	if b.Markers == nil {
		b.Markers = &Expression{}

		return b.Markers
	}

	m := &Expression{}
	//m.Parent = b.Markers.Tail
	//b.Markers.Tail.Children = append(b.Markers.Tail.Children, m)
	//b.Markers.Tail = m
	b.OrphanedChildren = append(b.OrphanedChildren, m)
	return m
}

func (b *Builder) AssignOrphanedChildren(m *Expression) {
	m.Children = append(m.Children, b.OrphanedChildren...)
	for _, c := range m.Children {
		c.Parent = m
	}
	b.OrphanedChildren = b.OrphanedChildren[:0]
	if m.IsDone {
		b.OrphanedChildren = append(b.OrphanedChildren, m)
	}
}

func (b *Builder) GetTokenType() state.ElementType {
	return b.CurrentToken.Type
}

func (b *Builder) AdvanceLexer() {
	t, _, v := b.Lexer.Lex()

	ty, ok := state.TokenTypes[t]
	if !ok {
		panic("wrong lek.Token")
	}

	b.CurrentToken = Token{ty, v}

}

func (b *Builder) Error(err string) {
	// TODO implement
}

type Parser struct {
}

func (p *Parser) ParseOrClause(b *Builder) bool {
	var marker *Expression
	if !p.ParseAndClause(b) {
		return false
	}

	for p.AdvanceIfMatches(b, state.OR_OPERATORS) {
		if marker == nil {
			marker = NewExpression()
			marker.Value = "OR"
			b.AssignOrphanedChildren(marker)
		}

		if !p.ParseAndClause(b) {
			// b.Errors probably need to panic or terminate parse
			b.Error("expected clause after OR keyword")
		}

		b.AssignOrphanedChildren(marker)
	}

	if marker != nil {
		marker.Done(state.OR_CLAUSE)
		b.AssignOrphanedChildren(marker)
	}

	return true
}

func (p *Parser) ParseAndClause(b *Builder) bool {
	var marker *Expression
	if !p.ParseTerminalClause(b) {
		return false
	}

	for p.AdvanceIfMatches(b, state.AND_OPERATORS) {
		if marker == nil {
			marker = NewExpression()
			marker.Value = "AND"
			b.AssignOrphanedChildren(marker)
		}

		if !p.ParseTerminalClause(b) {
			// b.Errors probably need to panic or terminate parse
			b.Error("expected clause after AND keyword")
			return false
		}
		b.AssignOrphanedChildren(marker)
	}

	if marker != nil {
		marker.Done(state.AND_CLAUSE)
		b.AssignOrphanedChildren(marker)
	}

	return true
}

func (p *Parser) ParseTerminalClause(b *Builder) bool {
	var marker *Expression
	if !p.ParseFieldName(b) {
		return false
	}

	ct := b.CurrentToken
	if p.AdvanceIfMatches(b, state.SIMPLE_OPERATORS) {
		marker = &Expression{}
		//marker.Precede(b)
		p.ParseOperand(b)
	}

	if marker != nil {
		marker.Value = ct.Value
		marker.Done(state.SIMPLE_CLAUSE)
		b.AssignOrphanedChildren(marker)
	}

	return true
}

func (p *Parser) ParseFieldName(b *Builder) bool {
	ct := b.CurrentToken
	if !p.AdvanceIfMatches(b, state.VALID_FIELD_NAMES) {
		b.Error("expected field name")
		return false
	}
	marker := b.Mark()
	marker.Value = ct.Value
	marker.Done(state.IDENTIFIER)
	return true
}

func (p *Parser) ParseOperand(b *Builder) bool {
	var marker *Expression
	parsed := true
	ct := b.CurrentToken
	if p.AdvanceIfMatches(b, state.LITERALS) {
		marker = b.Mark()
		marker.Value = ct.Value
		marker.Done(state.LITERAL)
	} else {
		//	marker.Drop()
		parsed = false
	}
	if !parsed {
		b.Error("expected either literal")
	}
	return parsed
}

func (p *Parser) AdvanceIfMatches(b *Builder, m map[state.ElementType]bool) bool {
	tt := b.GetTokenType()
	_, ok := m[tt]
	if ok {
		b.AdvanceLexer()
		return true
	}
	return false
}
