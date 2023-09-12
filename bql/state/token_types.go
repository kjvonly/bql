package state

type ElementType string

// TOKEN SETS

const STRING_LITERAL ElementType = "STRING_LITERAL"
const IDENTIFIER ElementType = "IDENTIFIER"

const EQ ElementType = "EQ"

var VALID_FIELD_NAMES = map[ElementType]bool{
	STRING_LITERAL: true,
	IDENTIFIER:     true,
}

var SIMPLE_OPERATORS = map[ElementType]bool{
	EQ: true,
}
