package state_test

import (
	"testing"

	ll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"launchpad.net/kjvonly-bql/bql/state"
)

func TestDoubleLinkedListCreated(t *testing.T) {
	//input := `book=(john, matthew) AND text=love OR text=world`

	tokens := ll.New()
	tokens.Add(state.Token{bqlIdentifier, "book"})

	// ident       "book"
	// eq          '='
	// lpar        '('
	// ident       "john"
	// comma       ','
	// ident       "matthew"
	// comma       ','
	// string      "and"
	// rpar        ')'
	// and         "AND"
	// ident       "text"
	// eq          '='
	// ident       "love"
	// or          "OR"
	// ident       "text"
	// eq          '='
	// ident       "world"

}
