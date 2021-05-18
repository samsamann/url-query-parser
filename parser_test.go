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
		{s: `filter[id][EQ]=24cb00fa-9f69-4107-852b-d7574dd217e4`},
		{s: `filter[number][GT]=123`},
		{s: `filter[a.b][EQ]=123`},
		{s: `filter[p.a.t.h][LIKE]=123`},
		{s: `sort=bar`},
		{s: `sort=-foo`},
		{s: `sort=foo,-bar`},
		{s: `page[offset]=2`},
		{s: `page[number]=1&page[size]=10`},
		{s: `field[type1]`},
		{s: `field[type1]=field1,field2`},
		{s: `include=foo`},
		{s: `include=foo,bar`},
		{s: `include=foo.bar`},
		{s: `include=first.foo.bar,second`},
	}

	for _, tt := range tests {
		p := NewParser(strings.NewReader(tt.s))
		_, err := p.Parse()
		if err != nil {
			t.Error(err)
		}
	}
}
