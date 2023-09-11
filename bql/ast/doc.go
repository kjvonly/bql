// ast package contains code to process the tokens into
// an abstract syntax tree used by the parser.

// @author Mikhail Golubev
//
// Slightly refactored JQL grammar. See original ANTLR parser grammar at:
// http://jira.stagingonserver.com/jira-project/jira-components/jira-core/src/main/antlr3/com/atlassian/jira/jql/parser/antlr/Jql.g
//
// query ::= or_clause [order_by]
// or_clause ::= and_clause {or_op and_clause}
// and_clause ::= not_expr {and_op not_expr}
// not_expr ::= not_op not_expr
//
//	| subclause
//	| terminal_clause
//
// subclause ::= "(" or_clause ")"
// terminal_clause ::= simple_clause
//
//	| was_clause
//	| changed_clause
//
// simple_clause ::= field simple_op value
// # although this is not mentioned in JQL manual, usage of both "from" and "to" predicates in "was" clause is legal
// was_clause ::= field "was" ["not"] ["in"] operand {history_predicate}
// changed_clause ::= field "changed" {history_predicate}
// simple_op ::= "="
//
//	| "!="
//	| "~"
//	| "!~"
//	| "<"
//	| ">"
//	| "<="
//	| ">="
//	| ["not"] "in"
//	| "is" ["not"]
//
// not_op ::= "not" | "!"
// and_op ::= "and" | "&&" | "&"
// or_op ::= "or" | "||" | "|"
// history_predicate ::= "from" operand
//
//	| "to" operand
//	| "by" operand
//	| "before" operand
//	| "after" operand
//	| "on" operand
//	| "during" operand
//
// field ::= string
//
//	| NUMBER
//	| CUSTOM_FIELD
//
// operand ::= empty
//
//	| string
//	| NUMBER
//	| func
//	| list
//
// empty ::= "empty" | "null"
// list ::= "(" operand {"," operand} ")"
// func ::= fname "(" arg_list ")"
// # function name can be even number (!)
// fname ::= string | NUMBER
// arg_list ::= argument {"," argument}
// argument ::= string | NUMBER
// string ::= SQUOTED_STRING
//
//	| QUOTED_STRING
//	| UNQOUTED_STRING
//
// order_by ::= "order" "by" sort_key {sort_key}
// sort_key ::= field ("asc" | "desc")
package ast
