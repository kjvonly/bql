package parser

import "launchpad.net/kjvonly-bql/bql/state"

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
