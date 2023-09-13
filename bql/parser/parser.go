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
	Next *Marker

	IsDropped bool
	IsDone    bool
	Type      state.ElementType
}

func (m *Marker) AddMarker(n *Marker) {
}

func (m *Marker) Drop() {
	m.IsDropped = true
}

func (m *Marker) Done(t state.ElementType) {
	{
		for n := m.Next; n != nil; n = n.Next {
			if !n.IsDone {
				//TODO should change panic to something else
				panic("all markers past this marker not done.")
			}
		}
	}
	m.IsDone = true
	m.Type = t
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
	b.Markers.Tail.Next = m
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

func (p *Parser) AdvanceIfMatches(b *Builder, m map[state.ElementType]bool) bool {
	_, ok := m[b.GetTokenType()]
	if ok {
		b.AdvanceLexer()
		return true
	}
	return false
}
