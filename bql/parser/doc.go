package parser

/**
 * 13f1ee1f
 * Parser implements ANTLR process in go
 * Used jetbrains implementation as inspiration.
 *
 * @author Mikhail Golubev
 *
 * Slightly refactored JQL grammar. See original ANTLR parser grammar at:
 * http://jira.stagingonserver.com/jira-project/jira-components/jira-core/src/main/antlr3/com/atlassian/jira/jql/parser/antlr/Jql.g
 *
 * query ::= or_clause [order_by]
 * or_clause ::= and_clause {or_op and_clause}
 * and_clause ::= not_expr {and_op not_expr}
 * not_expr ::= not_op not_expr
 *            | subclause
 *            | terminal_clause
 * subclause ::= "(" or_clause ")"
 * terminal_clause ::= simple_clause
 *                   | was_clause
 *                   | changed_clause
 * simple_clause ::= field simple_op value
 * # although this is not mentioned in JQL manual, usage of both "from" and "to" predicates in "was" clause is legal
 * was_clause ::= field "was" ["not"] ["in"] operand {history_predicate}
 * changed_clause ::= field "changed" {history_predicate}
 * simple_op ::= "="
 *             | "!="
 *             | "~"
 *             | "!~"
 *             | "<"
 *             | ">"
 *             | "<="
 *             | ">="
 *             | ["not"] "in"
 *             | "is" ["not"]
 * not_op ::= "not" | "!"
 * and_op ::= "and" | "&&" | "&"
 * or_op ::= "or" | "||" | "|"
 * history_predicate ::= "from" operand
 *                     | "to" operand
 *                     | "by" operand
 *                     | "before" operand
 *                     | "after" operand
 *                     | "on" operand
 *                     | "during" operand
 * field ::= string
 *         | NUMBER
 *         | CUSTOM_FIELD
 * operand ::= empty
 *           | string
 *           | NUMBER
 *           | func
 *           | list
 * empty ::= "empty" | "null"
 * list ::= "(" operand {"," operand} ")"
 * func ::= fname "(" arg_list ")"
 * # function name can be even number (!)
 * fname ::= string | NUMBER
 * arg_list ::= argument {"," argument}
 * argument ::= string | NUMBER
 * string ::= SQUOTED_STRING
 *          | QUOTED_STRING
 *          | UNQOUTED_STRING
 * order_by ::= "order" "by" sort_key {sort_key}
 * sort_key ::= field ("asc" | "desc")
 *
 */

/*
 * https://www.javacodegeeks.com/2017/09/guide-parsing-algorithms-terminology.html
 *
 *   Left-recursive Rules
 * In the context of parsers, an important feature is the support for left-recursive rules. This means that a rule starts with a reference to itself. Sometime this reference could also be indirect, that is to say it could appear in another rule referenced by the first one.
 *
 * Consider for example arithmetic operations. An addition could be described as two expression(s) separated by the plus (+) symbol, but the operands of the additions could be other additions.
 *
 * addition       : expression '+' expression
 * multiplication : expression '*' expression
 * // an expression could be an addition or a multiplication or a number
 * expression     : multiplication | addition | [0-9]+
 * In this example expression contains an indirect reference to itself via the rules addition and multiplication.
 *
 * This description also matches multiple additions like 5 + 4 + 3. That is because it can be interpreted as expression (5) ('+') expression(4+3) (the rule addition: the first expression corresponds to the option [0-9]+, the second one is another addition). And then 4 + 3 itself can be divided in its two components: expression(4) ('+') expression(3) (the rule addition:both the first and second expression corresponds to the option [0-9]+) .
 *
 * The problem is that left-recursive rules may not be used with some parser generators. The alternative is a long chain of expressions, that takes care also of the precedence of operators. A typical grammar for a parser that does not support such rules would look similar to this one:
 *
 *
 * expression     : addition
 * addition       : multiplication ('+' multiplication)*
 * multiplication : atom ('*' atom)*
 * atom           : [0-9]+
 * As you can see, the expressions are defined in the inverse order of precedence. So the parser would put the expression with the lower precedence at the lowest level of the three; thus they would be executed first.
 *
 * Some parser generators support direct left-recursive rules, but not indirect ones. Notice that usually the issue is with the parsing algorithm itself, that does not support left-recursive rules. So the parser generator may transform rules written as left-recursive in the proper way to make it work with its algorithm. In this sense, left-recursive support may be (very useful) syntactic sugar.
 *
 * How Left-recursive Rules Are Transformed
 * The specific way in which the rules are transformed vary from one parser generator to the other, however the logic remains the same. The expressions are divided in two groups: the ones with an operator and two operands and the atomic ones. In our example the only atomic expression is a number ([0-9]+), but it could also be an expression between parentheses ((5 + 4)). That is because in mathematics parentheses are used to increase the precedence of an expression.
 *
 * Once you have these two groups: you maintain the order of the members of the second group and reverse the order of the members of the first group. The reason is that humans reason on a first come, first serve basis: it is easier to write the expressions in their order of precedence.
 */
