package parser

import (
	"launchpad.net/kjvonly-bql/bql/state"
	"launchpad.net/kjvonly-bql/lex"
)

type Token struct {
	Token lex.Token
	Value string
}

type ElementType string

const (
	IDENTIFIER ElementType = "identifier"
)

type Element struct {
	Value interface{}
	Type  ElementType
}

func (e *Element) Add(Element) {

}

type Parser struct {
	Lexer   *lex.Lexer
	Element Element
}

func (p *Parser) Parse(input string) {

	p.Lexer = state.BQLLexer(input)
	e := Element

}

func (p *Parser) ParseQuery() {

}

func (p *Parser) advanceIfMatches() {
	t, _, v := p.Lexer.Lex()
	p.Element{Token{t, v}}
}
