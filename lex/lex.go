// Copyright 2017-2020 Denis Bernard <db047h@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package lex

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

// An EncodingError may be emitted by State.ReadRune upon reading invalid UTF-8 data.
type EncodingError struct {
	s string
}

func (e EncodingError) Error() string { return e.s }

// Encoding errors.
var (
	ErrNulChar     = &EncodingError{"invalid NUL character"}
	ErrInvalidRune = &EncodingError{"invalid UTF-8 encoding"}
	ErrInvalidBOM  = &EncodingError{"invalid BOM in the middle of the file"}
)

// ErrInvalidUnreadRune is returned by State.UnreadRune if the undo buffer is
// empty.
var ErrInvalidUnreadRune = errors.New("invalid use of UnreadRune")

// EOF is the return value from Next() when EOF is reached.
const EOF rune = -1

// Undo buffer constants.
const (
	BackupBufferSize = 16 // Size of the undo buffer.
	undoMask         = BackupBufferSize - 1
)

// A Token represents the type of a token. Custom lexers can use any value >= 0.
type Token int

// Error is the token type for error tokens.
const Error Token = -1

// queue is a FIFO queue.
type queue struct {
	items []item
	head  int
	tail  int
	count int
}

type item struct {
	t Token
	p int
	v interface{}
}

func (q *queue) push(t Token, p int, v interface{}) {
	if t == Error {
		if _, ok := v.(error); !ok {
			panic("token value must implement the error interface for Error tokens")
		}
	}
	if q.head == q.tail && q.count > 0 {
		items := make([]item, len(q.items)*2)
		copy(items, q.items[q.head:])
		copy(items[len(q.items)-q.head:], q.items[:q.head])
		q.head = 0
		q.tail = len(q.items)
		q.items = items
	}
	q.items[q.tail] = item{t, p, v}
	q.tail = (q.tail + 1) % len(q.items)
	q.count++
}

// pop pops the first item from the queue. Callers must check that q.count > 0 beforehand.
func (q *queue) pop() (Token, int, interface{}) {
	i := q.head
	q.head = (q.head + 1) % len(q.items)
	q.count--
	it := &q.items[i]
	return it.t, it.p, it.v
}

// Lexer wraps the public methods of a lexer. This interface is intended for
// parsers that call New(), then Lex() until EOF.
type Lexer state

// State holds the internal state of the lexer while processing a given input.
// Note that the public fields should only be accessed from custom StateFn
// functions.
type State state

type undo struct {
	p int
	r rune
	s int
}

type state struct {
	buf    [4 << 10]byte          // byte buffer
	undo   [BackupBufferSize]undo // undo buffer
	queue                         // Item queue
	f      *File
	line   int     // line count
	state  StateFn // current state
	init   StateFn // current initial-state function.
	offs   int     // offset of first byte in buffer
	r, w   int     // read/write indices
	ur, uh int     // undo buffer read pos and head
	ts     int     // token start offset
	ioErr  error   // if not nil, IO error @w
}

// A StateFn is a state function.
//
// If a StateFn returns nil, the lexer resets the current token starting offset
// then transitions back to its initial state function.
type StateFn func(l *State) StateFn

// NewLexer creates a new lexer associated with the given source file. A new
// lexer must be created for every source file to be lexed.
func NewLexer(f *File, init StateFn) *Lexer {
	s := &state{
		// initial q size must be an exponent of 2
		queue: queue{items: make([]item, 2)},
		f:     f,
		line:  1,
		init:  init,
		uh:    1,
	}

	// add line 1 to file
	f.AddLine(0, 1)
	// sentinel values
	for i := range s.undo {
		s.undo[i] = undo{-1, utf8.RuneSelf, 1}
	}

	return (*Lexer)(s)
}

// Init (re-)sets the initial state function for the lexer. It can be used by
// state functions to implement context switches (e.g. switch from accepting
// plain text to expressions in a template-like language). This function returns
// the previous initial state.
func (s *State) Init(initState StateFn) StateFn {
	prev := s.init
	s.init = initState
	return prev
}

// Lex reads source text and returns the next item until EOF.
//
// As a convention, once the end of file has been reached (or some
// non-recoverable error condition), Lex() must return a token type that
// indicates an EOF condition. A common strategy is to emit an Error token with
// io.EOF as a value.
func (l *Lexer) Lex() (Token, int, interface{}) {
	for l.count == 0 {
		st := (*State)(l)
		if l.state == nil {
			l.state = l.init(st)
		} else {
			l.state = l.state(st)
		}
	}
	return l.pop()
}

// File returns the File used as input for the lexer.
func (l *Lexer) File() *File {
	return l.f
}

// Emit emits a single token of the given type and value. offset is the file
// offset for the token (usually s.TokenPos()).
//
// If the emitted token is Error, the value must be an error interface.
func (s *State) Emit(offset int, t Token, value interface{}) {
	s.push(t, offset, value)
}

// Errorf emits an error token with type Error. The Item value is set to the
// result of calling fmt.Errorf(format, args...) and offset is the file offset.
func (s *State) Errorf(offset int, format string, args ...interface{}) {
	s.push(Error, offset, fmt.Errorf(format, args...))
}

// Next returns the next rune in the input stream. If the end of the input
// has ben reached it will return EOF. If an I/O error occurs other than io.EOF,
// it will report the I/O error by calling Errorf then return EOF.
//
// Next only returns valid runes or -1 to indicate EOF. It filters out invalid
// runes, nul bytes (0x00) and BOMs (U+FEFF) and reports them as errors by
// calling Errorf (except for a BOM at the beginning of the file which is simply
// ignored).
func (s *State) Next() rune {
	r, _, err := s.ReadRune()
	if err != nil {
		if err != io.EOF {
			s.Emit(s.Pos(), Error, err)
			s.ioErr = io.EOF
		}
		return EOF
	}
	return r
}

// ReadRune reads a single UTF-8 encoded Unicode character and returns the
// rune and its size in bytes. The returned rune is always valid. Invalid runes,
// NUL bytes and misplaced BOMs are filtered out and emitted as Error tokens.
//
// This function is here to facilitate interfacing with standard library
// scanners that need an io.RuneScanner. Custom lexers should use the Next
// function instead.
func (s *State) ReadRune() (rune, int, error) {
	// read from undo buffer
	if u := (s.ur + 1) & undoMask; u != s.uh {
		s.ur = u
		return s.undo[u].r, s.undo[u].s, nil
	}
again:
	for s.r+utf8.UTFMax > s.w && !utf8.FullRune(s.buf[s.r:s.w]) && s.ioErr == nil && s.w-s.r < len(s.buf) {
		s.fill()
	}

	off := s.offs + s.r

	// @ EOF
	if s.r == s.w {
		if s.Current() != EOF {
			s.pushUndo(off, EOF, 1)
		}
		return 0, 0, s.ioErr
	}

	// Common case: ASCII
	if b := s.buf[s.r]; b < utf8.RuneSelf {
		s.r++
		if b == 0 {
			s.Emit(off, Error, ErrNulChar)
			goto again
		}
		if b == '\n' {
			s.line++
			s.f.AddLine(off+1, s.line)
		}
		s.pushUndo(off, rune(b), 1)
		return rune(b), 1, nil
	}

	// UTF8
	r, w := utf8.DecodeRune(s.buf[s.r:s.w])
	s.r += w
	if r == utf8.RuneError && w == 1 {
		s.Emit(off, Error, ErrInvalidRune)
		goto again
	}

	// BOM only allowed as first rune in the file
	if r == 0xfeff {
		if off > 0 {
			s.Emit(off, Error, ErrInvalidBOM)
		}
		goto again
	}

	s.pushUndo(off, r, w)
	return r, w, nil
}

func (s *State) pushUndo(off int, r rune, sz int) {
	s.ur = s.uh
	s.undo[s.uh] = undo{off, r, sz}
	s.uh = (s.uh + 1) & undoMask
	s.undo[s.uh] = undo{-1, utf8.RuneSelf, 1}
}

// Backup reverts the last call to Next. Backup can be called at most
// (BackupBufferSize-1) times in a row (i.e. with no calls to Next in between).
// Calling Backup beyond the start of the undo buffer or at the beginning
// of the input stream will fail silently, Pos will return -1 (an invalid
// offset) and Current will return utf8.RuneSelf, a value impossible to get
// by any other means.
func (s *State) Backup() {
	if s.undo[s.ur].p == -1 {
		return
	}
	s.ur = (s.ur - 1) & undoMask
}

// UnreadRune reverts the last call to ReadRune. It is essentially the same as
// Backup except for the error return value.
func (s *State) UnreadRune() error {
	if s.undo[s.ur].p == -1 {
		return ErrInvalidUnreadRune
	}
	s.ur = (s.ur - 1) & undoMask
	return nil
}

// Current returns the last rune returned by State.Next.
func (s *State) Current() rune {
	return s.undo[s.ur].r
}

// Pos returns the byte offset of the last rune returned by State.Next.
// Returns -1 if no input has been read yet.
func (s *State) Pos() int {
	return s.undo[s.ur].p
}

func (s *State) fill() {
	// slide buffer contents
	if n := s.r; n > 0 {
		copy(s.buf[:], s.buf[n:s.w])
		s.offs += n
		s.w -= n
		s.r = 0
	}

	for i := 0; i < 100; i++ {
		n, err := s.f.Read(s.buf[s.w:len(s.buf)])
		s.w += n
		if err != nil {
			s.ioErr = err
			return
		}
		if s.w-s.r >= utf8.UTFMax {
			return
		}
	}

	s.ioErr = io.ErrNoProgress
}

// Peek returns the next rune in the input stream without consuming it. This
// is equivalent to calling Next followed by Backup. At EOF, it simply returns
// EOF.
func (s *State) Peek() rune {
	if s.Current() == EOF {
		return EOF
	}
	r := s.Next()
	s.Backup()
	return r
}

// StartToken sets offset as a token start offset. This is a utility function
// that when used in conjunction with TokenPos enables tracking of a token start
// position across a StateFn chain without having to manually keep track of it
// via closures or function parameters.
//
// This is typically called at the start of the initial state function, just
// after an initial call to next:
//
//	func stateInit(s *lexer.State) lexer.StateFn {
//		r := s.Next()
//		s.StartToken(s.Pos())
//		switch {
//		case r >= 'a' && r <= 'z':
//			return stateIdentifier
//		default:
//			// ...
//		}
//		return nil
//	}
//
//	func stateIdentifier(s *lexer.State) lexer.StateFn {
//		tokBuf := make([]rune, 0, 64)
//		for r := s.Current(); r >= 'a' && r <= 'z'; r = s.Next() {
//			tokBuf = append(tokBuf, r)
//		}
//		s.Backup()
//		// TokenPos gives us the token starting offset set in stateInit
//		s.Emit(s.TokenPos(), tokTypeIdentifier, string(tokBuf))
//		return nil
//	}
func (s *State) StartToken(offset int) {
	s.ts = offset
}

// TokenPos returns the last offset set by StartToken.
func (s *State) TokenPos() int {
	return s.ts
}
