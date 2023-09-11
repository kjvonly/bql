package parser

import (
	"launchpad.net/kjvonly-bql/lex"
)

type Token struct {
	Token lex.Token
	Value string
}

type Marker struct{}

type Build interface {
	Mark() Marker
	GetToken() ElementType
}

type Builder struct {
	Marks []Marker
}

// Mark adds a placeholder for new
func (b *Builder) Mark() Marker {
	m := Marker{}
	b.Marks = append(b.Marks, m)
	return m
}

func (b *Builder) GetToken() ElementType {
	return STRING_LITERAL
}

type Parser struct {
	Lexer *lex.Lexer
}

func (p *Parser) AdvanceIfMatches(b Builder, m map[ElementType]bool) bool {
	return false
}
