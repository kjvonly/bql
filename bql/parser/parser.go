package parser

import (
	"launchpad.net/kjvonly-bql/bql/state"
	"launchpad.net/kjvonly-bql/lex"
)

type Token struct {
	Type  state.ElementType
	Value interface{}
}

// NewMarker creates root Marker
func NewMarker() *Marker {
	return &Marker{}
}

func NewMarkerList() *MarkerList {
	return &MarkerList{}
}

type MarkerList struct {
	Head *Marker
	Tail *Marker
}

type Marker struct {
	Children []*Marker
	Parent   *Marker

	IsDropped bool
	IsDone    bool
	Type      state.ElementType
}

func (m *Marker) Drop() {
	m.IsDropped = true
}

func checkAllMarkersDoneOrDropped(n []*Marker) {
	for i := 0; i < len(n); i++ {
		if !n[i].IsDone && !n[i].IsDropped {
			//TODO should change panic to something else
			panic("all markers past this marker not done.")
		}

		checkAllMarkersDoneOrDropped(n[i].Children)
	}
}

func (m *Marker) Done(t state.ElementType) {
	checkAllMarkersDoneOrDropped(m.Children)
	m.IsDone = true
	m.Type = t
}

// Creates and returns a new marker starting immediately before the start of
// this marker and extending after its end. Can be called on a completed or
// a currently active marker.
//
// @return the new marker instance.
func (m *Marker) Precede(b *Builder) {
	b.Markers.Tail = m
}

type Build interface {
}

type Builder struct {
	Markers      *MarkerList
	Lexer        *lex.Lexer
	CurrentToken Token
}

func NewBuilder(lex *lex.Lexer) *Builder {
	return &Builder{
		Lexer:   lex,
		Markers: NewMarkerList(),
	}
}

// Mark adds a placeholder for new
func (b *Builder) Mark() *Marker {
	if b.Markers.Head == nil {
		b.Markers.Head = NewMarker()
		b.Markers.Tail = b.Markers.Head
		return b.Markers.Head
	}

	m := &Marker{}
	m.Parent = b.Markers.Tail
	b.Markers.Tail.Children = append(b.Markers.Tail.Children, m)
	b.Markers.Tail = m

	return m
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

func (p *Parser) ParseAndClause(b *Builder) bool {
	marker := b.Mark()
	if !p.ParseTerminalClause(b) {
		marker.Drop()
		return false
	}

	for p.AdvanceIfMatches(b, state.AND_OPERATORS) {
		marker.Precede(b)

		if !p.ParseTerminalClause(b) {
			// b.Errors probably need to panic or terminate parse
			b.Error("expected clause after AND keyword")
		}

	}

	if len(marker.Children) > 0 {
		marker.Done(state.AND_CLAUSE)
	} else {
		marker.Drop()
	}

	return true
}

func (p *Parser) ParseTerminalClause(b *Builder) bool {
	marker := b.Mark()
	if !p.ParseFieldName(b) {
		marker.Drop()
		return false
	}

	if p.AdvanceIfMatches(b, state.SIMPLE_OPERATORS) {
		marker.Precede(b)
		p.ParseOperand(b)
	}
	marker.Done(state.SIMPLE_CLAUSE)
	return true
}

func (p *Parser) ParseFieldName(b *Builder) bool {
	marker := b.Mark()
	if !p.AdvanceIfMatches(b, state.VALID_FIELD_NAMES) {
		b.Error("expected field name")
		marker.Drop()
		return false
	}
	marker.Done(state.IDENTIFIER)
	return true
}

func (p *Parser) ParseOperand(b *Builder) bool {
	marker := b.Mark()
	parsed := true
	if p.AdvanceIfMatches(b, state.LITERALS) {
		marker.Done(state.LITERAL)
	} else {
		marker.Drop()
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
