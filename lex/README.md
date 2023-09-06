# lex
Copied initial code from [here](https://github.com/db47h/lex) and modified.

## Overview

Package lex provides the core of a lexer built as a Deterministic Finite State
Automaton whose states and associated actions are implemented as functions.

Clients of the package only need to provide state functions specialized in
lexing the target language. The package provides facilities to stream input
from a io.Reader with up to 15 bytes look-ahead, as well as utility functions
commonly used in lexers.

The implementation is similar to https://golang.org/src/text/template/parse/lex.go.
See also Rob Pike's talk about combining states and actions into state
functions: https://talks.golang.org/2011/lex.slide.

## Release notes

### v1.2.1

Improvements to error handling:

`Lexer.Errorf` now generates `error` values instead of `string`.
`lexer.Emit` enforces the use of `error` values for `Error` tokens.

`Scanner` now implements the `io.RuneScanner` interface.

This is a minor API breakage that does not impact any known client code.

### v1.0.0

Initial release