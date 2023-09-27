package parser

import (
	"launchpad.net/kjvonly-bql/bql/state"
	"launchpad.net/kjvonly-bql/lex"
)

type Builder struct {
	Expression          *Expression
	Lexer               *lex.Lexer
	CurrentToken        Token
	OrphanedExpressions []*Expression
}

func NewBuilder(lex *lex.Lexer) *Builder {
	return &Builder{
		Lexer:      lex,
		Expression: &Expression{},
	}
}

func (b *Builder) AddExpression() *Expression {
	e := &Expression{}
	b.OrphanedExpressions = append(b.OrphanedExpressions, e)
	return e
}

func (b *Builder) AssignOrphanedExpressions(e *Expression) {
	e.Expressions = append(e.Expressions, b.OrphanedExpressions...)
	b.OrphanedExpressions = b.OrphanedExpressions[:0]
	if e.IsDone {
		b.OrphanedExpressions = append(b.OrphanedExpressions, e)
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
