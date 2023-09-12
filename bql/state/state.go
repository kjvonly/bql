package state

import (
	"strings"
	"unicode"

	"launchpad.net/kjvonly-bql/lex"
	"launchpad.net/kjvonly-bql/lex/state"
)

// Token types.
const (
	BqlEOF        lex.Token = iota // 0 EOF
	BqlSemiColon                   // 1 semi-colon, EOL
	BqlInt                         // 2 integer literal
	BqlFloat                       // 3 float literal
	BqlString                      // 4 quoted string
	BqlChar                        // 5 quoted char
	BqlIdentifier                  // 6 identifier
	BqlDot                         // 7 '.' field/method selector
	BqlRawChar                     // 8 any other single character
	BqlLPAR                        // 9 (
	BqlRPAR                        // 10 )
	BqlComma                       // 11 ,

	BqlEQ // 12 =

	BqlANDKeyword // 13 and
	BqlORKeyword  // 14 or
)

var TokenTypes = map[lex.Token]ElementType{
	lex.Error:     "error",
	BqlEOF:        "EOF",
	BqlSemiColon:  "semicolon",
	BqlInt:        "integer",
	BqlFloat:      "float",
	BqlString:     "STRING_LITERAL",
	BqlChar:       "char",
	BqlIdentifier: "ident",
	BqlDot:        "dot",
	BqlRawChar:    "raw char",
	BqlLPAR:       "LPAR",
	BqlRPAR:       "RPAR",
	BqlComma:      "comma",
	BqlEQ:         "eq",
	BqlANDKeyword: "and",
	BqlORKeyword:  "or",
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
	quotedString := state.QuotedString(BqlString)
	quotedChar := state.QuotedChar(BqlChar)
	ident := identifier()
	number := state.Number(BqlInt, BqlFloat, '.')

	return func(s *lex.State) lex.StateFn {
		// get current rune (read for us by the lexer upon entering the initial state)
		r := s.Next()
		pos := s.Pos()
		// THE big switch
		switch r {
		case lex.EOF:
			// End of file
			s.Emit(pos, BqlEOF, nil)
			return nil
		case '\n', ';':
			// transform EOLs to semi-colons
			s.Emit(pos, BqlSemiColon, ';')
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
			s.Emit(pos, BqlDot, nil)
			return nil

		// BQL
		case '(':
			s.Emit(pos, BqlLPAR, r)
			return nil
		case ')':
			s.Emit(pos, BqlRPAR, r)
			return nil

		case '=':
			s.Emit(pos, BqlEQ, r)
			return nil

		case ',':
			s.Emit(pos, BqlComma, r)
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
			s.Emit(pos, BqlRawChar, r)
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
			l.Emit(pos, BqlANDKeyword, string(b))
			return nil
		}

		if strings.ToLower(string(b)) == "or" {
			l.Emit(pos, BqlORKeyword, string(b))
			return nil
		}

		l.Emit(pos, BqlIdentifier, string(b))
		return nil
	}
}

// BQL: a lexer for a Bible Query Language language.
func BQLLexer(input string) *lex.Lexer {
	inputFile := lex.NewFile("example", strings.NewReader(input))
	return lex.NewLexer(inputFile, bqlInit())
}
