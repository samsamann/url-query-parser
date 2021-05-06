package querystring

import (
	"strings"
	"testing"
)

// Ensure the scanner can scan tokens correctly.
func TestScannerScan(t *testing.T) {
	var tests = []struct {
		s   string
		tok Token
		lit string
	}{
		// Special tokens (EOF, ILLEGAL, WS)
		{s: ``, tok: EOF},
		{s: `?`, tok: ILLEGAL, lit: `?`},
		{s: ` `, tok: WS, lit: " "},

		//operators
		{s: `EQ`, tok: EQUAL, lit: "EQ"},
	}

	for i, tt := range tests {
		s := NewScanner(strings.NewReader(tt.s))
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
