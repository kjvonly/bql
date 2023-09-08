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
