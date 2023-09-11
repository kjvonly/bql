package ast

import (
	"fmt"

	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"launchpad.net/kjvonly-bql/bql/state"
)

type Ast struct {
	currentPos int
	Tokens     *doublylinkedlist.List
	Elements   []interface{}
}

func (a *Ast) Generate(ll *doublylinkedlist.List) *Query {
	a.Tokens = ll
	q := Query{Expr: []Expr{}}
	a.ParseQuery()
	return &q
}

func (a *Ast) ParseQuery() error {
	a.currentPos = 0

	for a.currentPos != a.Tokens.Size() {
		tt, err := a.GetToken(a.currentPos)
		if err != nil {
			return err
		}
		_, so := state.SIMPLE_OPERATORS[tt.Token]
		if so {
			_, err := a.ParseOperator(tt)
			return err
		}
		a.currentPos = a.currentPos + 1
	}
	return nil
}

func (a *Ast) ParseOperator(tt state.Token) (bool, error) {
	if tt.Token == state.BqlEQ {
		return a.ParseEqStmt(tt)
	}
	return false, nil
}

func (a *Ast) ParseEqStmt(tt state.Token) (bool, error) {
	prevTokenIndex := a.currentPos - 1

	t, err := a.GetToken(prevTokenIndex)
	if err != nil {
		return false, err
	}

	if !state.IsStandardField(t.Value) {
		return false, fmt.Errorf("QUERY ERROR: = OPERATOR PRECEDED BY UNKNOWN FIELD %s", t.Value)
	}

	// link the IDENT FIELD with the EQ STMT
	a.Elements = append(a.Elements, EqualStmt{})
	return true, nil
}

func (a *Ast) GetToken(index int) (state.Token, error) {
	v, ok := a.Tokens.Get(index)
	if !ok {
		return state.Token{}, fmt.Errorf("INVALID INDEX AT '%d'", index)
	}
	return v.(state.Token), nil

}

// All node types implement the Node interface.
type Node interface {
	Print()
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
}

// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
	stmtNode()
}

type Ident struct {
	Node
	Name string
}

type EqualStmt struct {
	Field Ident
	Expr  Expr
}

type Query struct {
	Expr []Expr
}

// expr
func (*Query) exprNode() {}
func (*Ident) exprNode() {}

// stmt
func (x *EqualStmt) stmtNode() {}

// funcs for QueryExpr
func (q *Query) Print() {
	fmt.Printf("Query: \n")
	for _, exp := range q.Expr {
		exp.Print()
	}
}

// funcs for EqualStmt
func (e *EqualStmt) Print() {
	fmt.Printf("%s", "EqualStmt: ")
	e.Field.Print()
	e.Expr.Print()
	fmt.Printf("%s", "End of EqualStmt:")
}

// funcs for Idnet
func (id *Ident) IsField() bool { return state.IsStandardField(id.Name) }
func (id *Ident) Print()        { fmt.Printf("leaf: %s", id.Name) }
