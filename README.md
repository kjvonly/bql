# Bible Query Language (BQL)

BQL is a query language that makes searching through the bible easier and more efficient.

## Workflow

We use `git` for version control [master](https://code.launchpad.net/~man4christ/kjvonly-bql/+git/kjvonly-bql/+ref/master). However, since bql is a written in `go` launchpad only supports go imports if `breezy` is used. So releases will be tagged in `git` and added to the `breezy` project `trunk` with a tagged revision.

### Remote import paths

From he [docs](https://pkg.go.dev/cmd/go#hdr-Remote_import_paths)

```go
// Launchpad (Bazaar)
import "launchpad.net/project"
import "launchpad.net/project/series"
import "launchpad.net/project/series/sub/directory"

import "launchpad.net/~user/project/branch"
import "launchpad.net/~user/project/branch/sub/directory"
```

To reuse a package from `kjvonly-bql`

```go
import launchpad.net/kjvonly-bql/lex
```

## BQL Intro

A simple BQL query searching for all the verses containing `love` in the book of `John`

```sql
book="john" and text="love"
```

A query with a function that counts all the verses containing `love` in the book of `John` using the `count()` function.

```sql
count(book="john" and text="love")
```

## Code Structure

To write a query language one needs to be able to interpret, validate, and execute a query. This is accomplished in programming by tokenizing the text with a lexer, parsing the tokens with a Abstract Syntax Tree [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree), then walking the tree using the [visitor pattern](https://en.wikipedia.org/wiki/Visitor_pattern).

### [BQL LEXER](./lex)

The [lexer](./lex) tokenize the query into several predefined states. The states defined consist of [fields](#fields)

### MODELS

#### FIELDS

A field in BQL is a word that represents a `KJVonly` field.
