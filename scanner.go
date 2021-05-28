package querystring

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

var eof = rune(0)

type Scanner struct {
	// pos    Position
	reader *bufio.Reader
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(reader),
	}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (Token, string) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if isInteger(ch) {
		s.unread()
		return s.scanInteger()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '[':
		return BRACKET_OPEN, string(ch)
	case ']':
		return BRACKET_CLOSE, string(ch)
	case '=':
		return ASSIGN, string(ch)
	case '&':
		return AMPERSAND, string(ch)
	case '-':
		return DASH, string(ch)
	case ',':
		return COMMA, string(ch)
	case '.':
		return DOT, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isInteger(ch) && ch != '-' && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch strings.ToLower(buf.String()) {
	case keywords[FIELD]:
		return FIELD, buf.String()
	case keywords[FILTER]:
		return FILTER, buf.String()
	case keywords[SORT]:
		return SORT, buf.String()
	case keywords[PAGE]:
		return PAGE, buf.String()
	case keywords[INCLUDE]:
		return INCLUDE, buf.String()
	}

	switch buf.String() {
	case operators[EQUAL]:
		return EQUAL, buf.String()
	case operators[NEQUAL]:
		return NEQUAL, buf.String()
	case operators[LIKE]:
		return LIKE, buf.String()
	case operators[LT]:
		return LT, buf.String()
	case operators[LE]:
		return LE, buf.String()
	case operators[GT]:
		return GT, buf.String()
	case operators[GE]:
		return GE, buf.String()
	}

	switch buf.String() {
	case pageKeywords[OFFSET]:
		return OFFSET, buf.String()
	case pageKeywords[LIMIT]:
		return LIMIT, buf.String()
	case pageKeywords[NUMBER]:
		return NUMBER, buf.String()
	case pageKeywords[SIZE]:
		return SIZE, buf.String()
	}

	return IDENT, buf.String()
}

func (s *Scanner) scanInteger() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isInteger(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	return INT, buf.String()
}

func (s *Scanner) read() rune {
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { _ = s.reader.UnreadRune() }
