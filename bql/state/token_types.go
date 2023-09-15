package state

// TOKEN SETS

const STRING_LITERAL ElementType = "STRING_LITERAL"
const NUMBER_LITERAL ElementType = "NUMBER_LITERAL"

const IDENTIFIER ElementType = "IDENTIFIER"

// KEYWORDS

const AND_KEYWORD ElementType = "AND_KEYWORD"

// Operators
const EQ ElementType = "EQ"

var VALID_FIELD_NAMES = map[ElementType]bool{
	STRING_LITERAL: true,
	IDENTIFIER:     true,
}

var SIMPLE_OPERATORS = map[ElementType]bool{
	EQ: true,
}

var LITERALS = map[ElementType]bool{
	STRING_LITERAL: true,
	IDENTIFIER:     true,
	NUMBER_LITERAL: true,
}

var AND_OPERATORS = map[ElementType]bool{
	AND_KEYWORD: true,
}
