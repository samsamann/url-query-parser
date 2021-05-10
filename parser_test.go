package querystring

import (
	"strings"
	"testing"
)

func TestParserParse(t *testing.T) {
	var tests = []struct {
		s string
	}{
		{s: `filter[foo][EQ]=bar`},
		{s: `filter[number][EQ]=123`},
		{s: `sort=bar`},
		{s: `sort=-foo`},
		{s: `sort=foo,-bar`},
		{s: `page[offset]=2`},
		{s: `page[number]=1&page[size]=10`},
		{s: `field[type1]`},
		{s: `field[type1]=field1,field2`},
		{s: `include=foo,bar`},
	}

	for _, tt := range tests {
		p := NewParser(strings.NewReader(tt.s))
		_, err := p.Parse()
		if err != nil {
			t.Error(err)
		}
	}
}
