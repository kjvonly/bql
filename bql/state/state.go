package state

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"launchpad.net/kjvonly-bql/lex"
	"launchpad.net/kjvonly-bql/lex/state"
)

// Token types.
const (
	bqlEOF        lex.Token = iota // 0 EOF
	bqlSemiColon                   // 1 semi-colon, EOL
	bqlInt                         // 2 integer literal
	bqlFloat                       // 3 float literal
	bqlString                      // 4 quoted string
	bqlChar                        // 5 quoted char
	bqlIdentifier                  // 6 identifier
	bqlDot                         // 7 '.' field/method selector
	bqlRawChar                     // 8 any other single character
	bqlLPAR                        // 9 (
	bqlRPAR                        // 10 )
	bqlComma                       // 11 ,

	bqlEQ // 12 =

	bqlANDKeyword // 13 and
	bqlORKeyword  // 14 or
)

var tokNames = map[lex.Token]string{
	lex.Error:     "error",
	bqlEOF:        "EOF",
	bqlSemiColon:  "semicolon",
	bqlInt:        "integer",
	bqlFloat:      "float",
	bqlString:     "string",
	bqlChar:       "char",
	bqlIdentifier: "ident",
	bqlDot:        "dot",
	bqlRawChar:    "raw char",
	bqlLPAR:       "lpar",
	bqlRPAR:       "rpar",
	bqlComma:      "comma",
	bqlEQ:         "eq",
	bqlANDKeyword: "and",
	bqlORKeyword:  "or",
}

// bqlInit returns the initial state function for our language.
// We implement it as a closure so that we can initialize state functions from
// the state package and take advantage of buffer pre-allocation.
func bqlInit() lex.StateFn {
	// Note that because of the buffer pre-allocation mentioned above, reusing
	// any of these variables in multiple goroutines is not safe. i.e. do not
	// turn these into global variables.
	// Instead, call tgInit() to get a new initial state function for each lexer
	// running in a goroutine.
	quotedString := state.QuotedString(bqlString)
	quotedChar := state.QuotedChar(bqlChar)
	ident := identifier()
	number := state.Number(bqlInt, bqlFloat, '.')

	return func(s *lex.State) lex.StateFn {
		// get current rune (read for us by the lexer upon entering the initial state)
		r := s.Next()
		pos := s.Pos()
		// THE big switch
		switch r {
		case lex.EOF:
			// End of file
			s.Emit(pos, bqlEOF, nil)
			return nil
		case '\n', ';':
			// transform EOLs to semi-colons
			s.Emit(pos, bqlSemiColon, ';')
			return nil
		case '"':
			return quotedString
		case '\'':
			return quotedChar
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return number
		case '.':
			// we want to distinguish a float starting with a leading dot from a dot used as
			// a field/method selector between two identifiers.
			if r = s.Peek(); r >= '0' && r <= '9' {
				// dot followed by a digit => floating point number
				return number
			}
			// for a dot followed by any other interesting char, we emit it as-is
			s.Emit(pos, bqlDot, nil)
			return nil

		// BQL
		case '(':
			s.Emit(pos, bqlLPAR, r)
			return nil
		case ')':
			s.Emit(pos, bqlRPAR, r)
			return nil

		case '=':
			s.Emit(pos, bqlEQ, r)
			return nil

		case ',':
			s.Emit(pos, bqlComma, r)
			return nil
		}

		// we're left with identifiers, spaces and raw chars.
		switch {
		case unicode.IsSpace(r):
			// eat spaces
			for r = s.Next(); unicode.IsSpace(r); r = s.Next() {
			}
			s.Backup()
			return nil
		case unicode.IsLetter(r) || r == '_':
			// r starts an identifier
			return ident
		default:
			s.Emit(pos, bqlRawChar, r)
			return nil
		}
	}
}

func identifier() lex.StateFn {
	// preallocate a buffer to store the identifier. It will end-up being at
	// least as large as the largest identifier scanned.
	b := make([]rune, 0, 64)
	return func(l *lex.State) lex.StateFn {
		pos := l.Pos()
		// reset buffer and add first char
		b = append(b[:0], l.Current())
		// read identifier
		for r := l.Next(); unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'; r = l.Next() {
			b = append(b, r)
		}

		// the character returned by the last call to next is not part of the identifier. Undo it.
		l.Backup()

		if strings.ToLower(string(b)) == "and" {
			l.Emit(pos, bqlANDKeyword, string(b))
			return nil
		}

		if strings.ToLower(string(b)) == "or" {
			l.Emit(pos, bqlORKeyword, string(b))
			return nil
		}

		l.Emit(pos, bqlIdentifier, string(b))
		return nil
	}
}

type Token struct {
	Token lex.Token
	Value string
}

// BQL: a lexer for a Bible Query Language language.
func BQLLexer(input string) map[string]string {
	// initialize lex.
	inputFile := lex.NewFile("example", strings.NewReader(input))
	l := lex.NewLexer(inputFile, bqlInit())

	// loop over each token
	for tt, _, v := l.Lex(); tt != bqlEOF; tt, _, v = l.Lex() {
		// print the token type and value.
		switch v := v.(type) {
		case nil:
			fmt.Println(tokNames[tt])
		case string:
			fmt.Printf("%-12s%s\n", tokNames[tt], strconv.Quote(v))
		case rune:
			fmt.Printf("%-12s%s\n", tokNames[tt], strconv.QuoteRune(v))
		default:
			fmt.Printf("%-12s%s\n", tokNames[tt], v)
		}
	}

	return nil
}
