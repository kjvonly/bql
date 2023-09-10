package state_test

import (
	"testing"

	ll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"launchpad.net/kjvonly-bql/bql/state"
)

func TestDoubleLinkedListCreated(t *testing.T) {
	input := `book = john and text = love`
	expectedLl := ll.New()
	expectedLl.Add(state.Token{state.BqlIdentifier, "book"})
	expectedLl.Add(state.Token{state.BqlEQ, "="})
	expectedLl.Add(state.Token{state.BqlIdentifier, "john"})
	expectedLl.Add(state.Token{state.BqlANDKeyword, "and"})
	expectedLl.Add(state.Token{state.BqlIdentifier, "text"})
	expectedLl.Add(state.Token{state.BqlEQ, "="})
	expectedLl.Add(state.Token{state.BqlIdentifier, "love"})

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

		if ev.Token != av.Token || ev.Value != av.Value {
			t.Fatalf("linked lists do not match. expected %+v, got: %+v", ev, av)
		}
	}
}
