package querystring

type Token int

const (
	EOF Token = iota
	ILLEGAL
	WS
	IDENT
	INT

	AMPERSAND
	ASSIGN
	BRACKET_OPEN
	BRACKET_CLOSE
	COMMA
	DASH

	// keywords
	FIELD
	FILTER
	SORT
	PAGE
	INCLUDE

	// page keywords
	OFFSET
	LIMIT
	NUMBER
	SIZE

	// ops
	EQUAL
	NEQUAL
	LIKE
)

var general = map[Token]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	WS:      " ",
	IDENT:   "IDENT",
	INT:     "INTEGER",

	AMPERSAND:     "&",
	ASSIGN:        "=",
	BRACKET_OPEN:  "[",
	BRACKET_CLOSE: "]",
	COMMA:         ",",
	DASH:          "-",
}

var keywords = map[Token]string{
	FIELD:   "field",
	FILTER:  "filter",
	SORT:    "sort",
	PAGE:    "page",
	INCLUDE: "include",
}

var pageKeywords = map[Token]string{
	OFFSET: "offset",
	LIMIT:  "limit",
	NUMBER: "number",
	SIZE:   "size",
}

var operators = map[Token]string{
	EQUAL:  "EQ",
	NEQUAL: "NEQ",
	LIKE:   "LIKE",
}
