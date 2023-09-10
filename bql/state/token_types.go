package state

import "launchpad.net/kjvonly-bql/lex"

var SIMPLE_OPERATORS = map[lex.Token]bool{
	BqlEQ: true,
}
