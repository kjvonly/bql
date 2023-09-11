package parser

type ElementType string

// TOKEN SETS

// TODO switch Idnetifier name to STRING_LITERAL
const STRING_LITERAL ElementType = "STRING_LITERAL"

var VALID_FIELD_NAMES = map[ElementType]bool{
	STRING_LITERAL: true,
}
