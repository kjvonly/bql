package parser

import (
	"launchpad.net/kjvonly-bql/bql/state"
	"launchpad.net/kjvonly-bql/lex"
)

type Token struct {
	Type  state.ElementType
	Value interface{}
}

type Marker struct{}

type Build interface {
	Mark() Marker
	GetTokenType() state.ElementType
	AdvanceLexer()
}

type Builder struct {
	Marks        []Marker
	Lexer        *lex.Lexer
	CurrentToken Token
}

func NewBuilder(lex *lex.Lexer) *Builder {
	return &Builder{Lexer: lex}
}

// Mark adds a placeholder for new
func (b *Builder) Mark() Marker {
	m := Marker{}
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

type Parser struct {
}

func (p *Parser) AdvanceIfMatches(b Builder, m map[state.ElementType]bool) bool {
	_, ok := m[b.GetTokenType()]
	if ok {
		b.AdvanceLexer()
		return true
	}
	return false
}
