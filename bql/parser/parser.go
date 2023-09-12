package parser

import (
	"launchpad.net/kjvonly-bql/bql/state"
	"launchpad.net/kjvonly-bql/lex"
)

type Token struct {
	Type  state.ElementType
	Value interface{}
}

type Marker struct {
	Dropped bool
	Type    state.ElementType
}

func (m *Marker) Drop() {
	m.Dropped = true
}

func (m *Marker) Done(t state.ElementType) {
	// TODO check everything after this is done
	m.Type = t
}

type Build interface {
}

type Builder struct {
	Marks        []*Marker
	Lexer        *lex.Lexer
	CurrentToken Token
}

func NewBuilder(lex *lex.Lexer) *Builder {
	return &Builder{Lexer: lex}
}

// Mark adds a placeholder for new
func (b *Builder) Mark() *Marker {
	m := &Marker{}
	b.Marks = append(b.Marks, m)
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
