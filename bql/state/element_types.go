package state

type ElementType string

const SIMPLE_CLAUSE ElementType = "SIMPLE_CLAUSE"
const AND_CLAUSE ElementType = "AND_CLAUSE"
const OR_CLAUSE ElementType = "OR_CLAUSE"

const QUERY ElementType = "QUERY"
const LITERAL ElementType = "LITERAL"
