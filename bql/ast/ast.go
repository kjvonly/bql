package ast

// All node types implement the Node interface.
type Node interface {
	Position() int
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
	Name string
}

type EqualStmt struct {
	Field Ident
	Expr  Expr
}

type ArrayType struct{}

// expr
func (*Ident) exprNode()     {}
func (*ArrayType) exprNode() {}

// stmt
func (x *EqualStmt) stmtNode() {}
