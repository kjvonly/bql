package parser

import (
	"launchpad.net/kjvonly-bql/bql/state"
	"launchpad.net/kjvonly-bql/lex"
)

type Token struct {
	Type  state.ElementType
	Value interface{}
}

type Expression struct {
	Expressions []*Expression
	Parent      *Expression
	IsDone      bool
	Type        state.ElementType
	Value       interface{}
}

func checkAllExpressionsDone(es []*Expression) {
	for i := 0; i < len(es); i++ {
		if !es[i].IsDone {
			//TODO should change panic to something else
			panic("all markers past this marker not done.")
		}

		checkAllExpressionsDone(es[i].Expressions)
	}
}

func (e *Expression) Done(t state.ElementType) {
	checkAllExpressionsDone(e.Expressions)
	e.IsDone = true
	e.Type = t
}

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
	for _, c := range e.Expressions {
		c.Parent = e
	}
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

type Parser struct{}

func (p *Parser) ParseOrClause(b *Builder) bool {
	var e *Expression
	if !p.ParseAndClause(b) {
		return false
	}

	for p.AdvanceIfMatches(b, state.OR_OPERATORS) {
		if e == nil {
			e = &Expression{}
			e.Value = "OR"
			b.AssignOrphanedExpressions(e)
		}

		if !p.ParseAndClause(b) {
			// b.Errors probably need to panic or terminate parse
			b.Error("expected clause after OR keyword")
		}

		b.AssignOrphanedExpressions(e)
	}

	if e != nil {
		e.Done(state.OR_CLAUSE)
		b.AssignOrphanedExpressions(e)
	}

	return true
}

func (p *Parser) ParseAndClause(b *Builder) bool {
	var e *Expression
	if !p.ParseTerminalClause(b) {
		return false
	}

	for p.AdvanceIfMatches(b, state.AND_OPERATORS) {
		if e == nil {
			e = &Expression{}
			e.Value = "AND"
			b.AssignOrphanedExpressions(e)
		}

		if !p.ParseTerminalClause(b) {
			// TODO b.Errors probably need to panic or terminate parse
			b.Error("expected clause after AND keyword")
			return false
		}
		b.AssignOrphanedExpressions(e)
	}

	if e != nil {
		e.Done(state.AND_CLAUSE)
		b.AssignOrphanedExpressions(e)
	}

	return true
}

func (p *Parser) ParseTerminalClause(b *Builder) bool {
	var e *Expression
	if !p.ParseFieldName(b) {
		return false
	}

	ct := b.CurrentToken
	if p.AdvanceIfMatches(b, state.SIMPLE_OPERATORS) {
		e = &Expression{}
		p.ParseOperand(b)
	}

	if e != nil {
		e.Value = ct.Value
		e.Done(state.SIMPLE_CLAUSE)
		b.AssignOrphanedExpressions(e)
	}

	return true
}

func (p *Parser) ParseFieldName(b *Builder) bool {
	ct := b.CurrentToken
	if !p.AdvanceIfMatches(b, state.VALID_FIELD_NAMES) {
		b.Error("expected field name")
		return false
	}
	e := b.AddExpression()
	e.Value = ct.Value
	e.Done(state.IDENTIFIER)
	return true
}

func (p *Parser) ParseOperand(b *Builder) bool {
	var e *Expression
	parsed := true
	ct := b.CurrentToken
	if p.AdvanceIfMatches(b, state.LITERALS) {
		e = b.AddExpression()
		e.Value = ct.Value
		e.Done(state.LITERAL)
	} else {
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
