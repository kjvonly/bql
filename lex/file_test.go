package lex_test

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/width"
	"launchpad.net/kjvonly-bql/lex"
)

// This example shows how one could use File.GetLineBytes to display nicely
// formatted error messages.
// For the example's sake, we use a dummy lexer that errors on digits, newlines
// and EOF.
func ExampleFile_GetLineBytes() {
	const (
		tokEOF = iota
		tokAny
	)
	expectLT := func(s *lex.State) lex.StateFn {
		// digits are followed by a < in order to test proper Seek operation in input.
		if s.Next() != '<' {
			panic("seek to original pos failed")
		}
		return nil
	}
	input := "＃〄 - Hello 世界 1<\ndéjà vu 2<"
	f := lex.NewFile("INPUT", strings.NewReader(input))
	l := lex.NewLexer(f, func(s *lex.State) lex.StateFn {
		switch r := s.Next(); r {
		case lex.EOF:
			s.Errorf(s.Pos(), "some error @EOF")
			s.Emit(s.Pos(), tokEOF, nil)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s.Errorf(s.Pos(), "digit")
			return expectLT
		case '\n':
			s.Errorf(s.Pos(), "newline")
		default:
			s.Emit(s.Pos(), tokAny, r)
		}
		return nil
	})
	for {
		tok, p, v := l.Lex()
		if tok == tokEOF {
			break
		}
		if tok == lex.Error {
			reportError(f, p, v.(error).Error())
		}
	}

	// The following output will display correctly only with monospaced fonts
	// and a UTF-8 locale. The caret alignment will also be off with some fonts
	// like Fira Code and East Asian characters.

	// Output:
	// INPUT:1:23: error digit
	// |＃〄 - Hello 世界 1<
	// |                  ^
	// INPUT:1:25: error newline
	// |＃〄 - Hello 世界 1<
	// |                    ^
	// INPUT:2:11: error digit
	// |déjà vu 2<
	// |        ^
	// INPUT:2:13: error some error @EOF
	// |déjà vu 2<
	// |          ^
}

// reportError reports a lexing error in the form:
//
//	file:line:col: error description
//		source line where the error occurred followed by a line with a carret at the position of the error.
//						      ^
func reportError(f *lex.File, p int, msg string) {
	pos := f.Position(p)
	fmt.Printf("%s: error %s\n", pos, msg)
	l, err := f.GetLineBytes(p)
	if err != nil {
		return
	}
	b := pos.Column - 1
	if b > len(l) {
		b = len(l)
	}
	fmt.Printf("|%s\n", l)
	fmt.Printf("|%*c^\n", getWidth(l[:b]), ' ')
	// or make it red!
	// fmt.Printf("|%*c\x1b[31m^\x1b[0m\n", getWidth(l[:b]), ' ')
}

// getWidth computes the width in text cells of a given byte slice.
// (supposing rendering with a UTF-8 locale and monospaced font)
func getWidth(l []byte) int {
	w := 0
	for i := 0; i < len(l); {
		r, s := utf8.DecodeRune(l[i:])
		i += s
		if !unicode.IsGraphic(r) {
			continue
		}
		p := width.LookupRune(r)
		switch p.Kind() {
		case width.EastAsianFullwidth, width.EastAsianWide:
			w += 2
		case width.EastAsianAmbiguous:
			w += 1 // depends on user locale. 2 if locale is CJK, 1 otherwise.
		default:
			w += 1
		}
	}
	return w
}
