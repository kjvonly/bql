package state

type StandardField string

const (
	BOOK StandardField = "BOOK"
	TEXT               = "TEXT"
)

var STANDARD_FIELDS map[string]StandardField = map[string]StandardField{
	"book": BOOK,
}

func IsStandardField(lookup string) bool {
	_, ok := STANDARD_FIELDS[lookup]
	return ok
}
