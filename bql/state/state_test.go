package state_test

import (
	"testing"

	ll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"launchpad.net/kjvonly-bql/bql/state"
)

func TestDoubleLinkedListCreated(t *testing.T) {
	//input := `book=(john, matthew) AND text=love OR text=world`
	input := `book`
	expectedLl := ll.New()
	expectedLl.Add(state.Token{state.BqlIdentifier, "book"})

	actualLl := state.BQLLexer(input)
	expectedIter := expectedLl.Iterator()
	for expectedIter.Next() {
		index := expectedIter.Index()
		ev := expectedIter.Value().(state.Token)
		actualValue, found := actualLl.Get(index)

		if !found {
			t.Fatalf("linked lists do not match. expected %+v, got: empty", ev)
		}
		var av state.Token = actualValue.(state.Token)

		if !found || ev.Token != av.Token {
			t.Fatalf("linked lists do not match. expected %+v, got: %+v", ev, av)
		}
	}
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
