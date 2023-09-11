package parser_test

import (
	"testing"

	"launchpad.net/kjvonly-bql/bql/parser"
)

func TestParseFieldName(t *testing.T) {

	b := parser.Builder{}
	_ = b.Mark()

	if len(b.Marks) != 1 {
		t.Fatalf("Should have 1 mark but has %d marks", len(b.Marks))
	}
}
