package state_test

import (
	"testing"

	ll "github.com/emirpasic/gods/lists/doublylinkedlist"
)

func TestHellow(t *testing.T) {
	//input := `book=(john, matthew) AND text=love OR text=world`

	ll.New()
	// expextedMap["ident   "] =  "\"book\""
	// expextedMap["raw char"] =  "'='"
	// expextedMap["raw char"] =  "'('"
	// expextedMap["ident   "] =  "\"john\""
	// expextedMap["raw char"] =  "','"
	// expextedMap["ident   "] =  "\"matthew\""
	// expextedMap["raw char"] =  "')'"
	// expextedMap["ident   "] =  "\"text\""
	// expextedMap["raw char"] =  "'='"
	// expextedMap["ident   "] =  "\"love\""
	// expextedMap["ident   "] =  "\"OR\""
	// expextedMap["ident   "] =  "\"text\""
	// expextedMap["raw char"] =  "'='"
	// expextedMap["ident   "] =  "\"world\""

}
