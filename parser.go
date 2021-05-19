package querystring

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	DEFAULT_PAGE_SIZE uint = 20
	MAX_PAGE_SIZE     uint = 100
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok  Token  // last read token
		lit  string // last read literal
		used bool   // buffer used
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Parse() (*QuerySpec, error) {
	return p.ParseWithPS(DEFAULT_PAGE_SIZE, MAX_PAGE_SIZE)
}

// Parse parses a query string
func (p *Parser) ParseWithPS(defaultPageSize, maxPageSize uint) (*QuerySpec, error) {
	querySpec := NewQuerySpec(defaultPageSize, maxPageSize)
	tok, lit := p.scanIgnoreWhitespace()
	for tok != EOF {
		switch tok {
		case FIELD:
			err := p.parseFieldSpec(querySpec.Field)
			if err != nil {
				return nil, err
			}
		case FILTER:
			filter, err := p.parseFilterSpec()
			if err != nil {
				return nil, err
			}
			querySpec.Filter = append(querySpec.Filter, filter)
		case SORT:
			sortSpecs, err := p.parseSortSpec()
			if err != nil {
				return nil, err
			}
			querySpec.Sort = append(querySpec.Sort, sortSpecs...)
		case PAGE:
			err := p.parsePageSpec(querySpec.PageSpec)
			if err != nil {
				return nil, err
			}
		case INCLUDE:
			sortSpecs, err := p.parseIncludeSpec()
			if err != nil {
				return nil, err
			}
			querySpec.IncludeSpec = append(querySpec.IncludeSpec, sortSpecs...)
		default:
			return nil, tokenError(lit, tokenMapToSlice(keywords)...)
		}

		tok, _ = p.scan()
		if tok != AMPERSAND {
			break
		}
		tok, lit = p.scan()
	}

	return querySpec, nil
}

func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.used {
		p.buf.used = false
		return p.buf.tok, p.buf.lit
	}

	tok, lit = p.s.Scan()
	p.buf.tok, p.buf.lit = tok, lit

	return
}

func (p *Parser) parseFilterSpec() (*FilterSpec, error) {
	filterSpec := new(FilterSpec)
	if tok, lit := p.scan(); tok != BRACKET_OPEN {
		return nil, tokenError(lit, general[BRACKET_OPEN])
	}

	buildPath := func(parent *pathElement, seg string) *pathElement {
		newSeg := &pathElement{segemnt: seg}
		parent.child = newSeg
		return newSeg
	}

	tok, lit := p.scan()
	if tok == IDENT || isPageKeyword(tok) {
		var rootElement pathElement = pathElement{segemnt: "."}
		var nextSeq = buildPath(&rootElement, lit)
		tok, lit = p.scan()
		for {
			if tok == IDENT || isPageKeyword(tok) {
				nextSeq = buildPath(nextSeq, lit)
				tok, lit = p.scan()
			} else if tok == DOT {
				tok, lit = p.scan()
			} else {
				p.unscan()
				break
			}
		}
		filterSpec.Field = rootElement
	} else {
		return nil, tokenError(lit, "field name", "field path")
	}

	if tok, lit := p.scan(); tok != BRACKET_CLOSE {
		return nil, tokenError(lit, general[BRACKET_CLOSE])
	}

	if tok, lit = p.scan(); tok != BRACKET_OPEN {
		return nil, tokenError(lit, general[BRACKET_OPEN])
	}

	tok, lit = p.scan()
	if !isOperator(tok) {
		return nil, tokenError(lit, tokenMapToSlice(operators)...)
	}
	filterSpec.Operator = lit

	if tok, lit := p.scan(); tok != BRACKET_CLOSE {
		return nil, tokenError(lit, general[BRACKET_CLOSE])
	}

	if tok, lit := p.scan(); tok != ASSIGN {
		return nil, tokenError(lit, general[ASSIGN])
	}

	tok, lit = p.scan()
	if tok != IDENT && tok != INT && !isOperator(tok) && !isPageKeyword(tok) {
		return nil, tokenError(lit, "value")
	}
	val := lit
	for {
		tok, lit = p.scan()
		if tok == AMPERSAND || tok == EOF || tok == ILLEGAL {
			p.unscan()
			break
		}
		val += lit
	}
	filterSpec.Value = val

	return filterSpec, nil
}

func (p *Parser) parseSortSpec() ([]*SortSpec, error) {
	if tok, lit := p.scan(); tok != ASSIGN {
		return nil, tokenError(lit, general[ASSIGN])
	}

	sortSpecs := make([]*SortSpec, 0)
	for {
		sort := new(SortSpec)
		tok, lit := p.scan()
		if tok == DASH {
			sort.Desc = true
			tok, lit = p.scan()
		}

		if tok == IDENT {
			sort.Field = lit
		} else {
			return nil, tokenError(lit, general[DASH], general[IDENT])
		}
		sortSpecs = append(sortSpecs, sort)
		if tok, _ = p.scan(); tok != COMMA {
			p.unscan()
			break
		}
	}
	return sortSpecs, nil
}

func (p *Parser) parsePageSpec(pageSpec *PageSpec) error {
	if tok, lit := p.scan(); tok != BRACKET_OPEN {
		return tokenError(lit, general[BRACKET_OPEN])
	}
	tok, lit := p.scan()
	if isPageKeyword(tok) {
		pageParam := tok
		if tok, lit = p.scan(); tok != BRACKET_CLOSE {
			return tokenError(lit, general[BRACKET_CLOSE])
		}
		if tok, lit = p.scan(); tok != ASSIGN {
			return tokenError(lit, general[ASSIGN])
		}
		if tok, lit = p.scan(); tok != INT {
			return tokenError(lit, general[INT])
		} else {
			i, _ := strconv.Atoi(lit)
			switch pageParam {
			case OFFSET:
				pageSpec.Offset = i
			case LIMIT:
				pageSpec.Limit = i
			case NUMBER:
				pageSpec.Number = i
			case SIZE:
				pageSpec.Size = i
			}
		}
		return nil
	}

	return tokenError(lit, tokenMapToSlice(pageKeywords)...)
}

func (p *Parser) parseFieldSpec(fieldSpec FieldSpec) error {
	if tok, lit := p.scan(); tok != BRACKET_OPEN {
		return tokenError(lit, general[BRACKET_OPEN])
	}
	tok, lit := p.scan()
	if tok != IDENT && !isOperator(tok) && !isPageKeyword(tok) {
		return tokenError(lit, "type")
	}
	entityType := lit
	if tok, lit = p.scan(); tok != BRACKET_CLOSE {
		return tokenError(lit, general[BRACKET_CLOSE])
	}
	fieldSpec[entityType] = make([]string, 0)
	tok, lit = p.scan()
	if tok == AMPERSAND {
		p.unscan()
		return nil
	} else if tok == EOF {
		return nil
	} else if tok != ASSIGN {
		return tokenError(lit, general[ASSIGN])
	}

	for {
		var field string
		tok, lit := p.scan()
		if tok == IDENT {
			field = lit
		} else {
			return tokenError(lit, general[DASH], general[IDENT])
		}
		fieldSpec[entityType] = append(fieldSpec[entityType], field)
		if tok, _ = p.scan(); tok != COMMA {
			p.unscan()
			break
		}
	}
	return nil
}

func (p *Parser) parseIncludeSpec() (IncludeSpec, error) {
	if tok, lit := p.scan(); tok != ASSIGN {
		return nil, tokenError(lit, general[ASSIGN])
	}

	includeSpecs := make(IncludeSpec, 0)

	buildPath := func(parent *pathElement, seg string) *pathElement {
		newSeg := &pathElement{segemnt: seg}
		parent.child = newSeg
		return newSeg
	}

	tok, lit := p.scan()
	if tok == IDENT || isPageKeyword(tok) {
		for {
			var rootElement pathElement = pathElement{segemnt: "."}
			var nextSeq = buildPath(&rootElement, lit)
			tok, lit = p.scan()
			for {
				if tok == IDENT || isPageKeyword(tok) {
					nextSeq = buildPath(nextSeq, lit)
					tok, lit = p.scan()
				} else if tok == DOT {
					tok, lit = p.scan()
				} else {
					break
				}
			}
			includeSpecs = append(includeSpecs, rootElement)
			if tok == COMMA {
				_, lit = p.scan()
				continue
			}
			p.unscan()
			break
		}
	} else {
		return nil, tokenError(lit, general[IDENT])
	}
	return includeSpecs, nil
}

func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	for tok == WS {
		tok, lit = p.scan()
	}
	return
}

func (p *Parser) unscan() { p.buf.used = true }

func tokenMapToSlice(tok map[Token]string) []string {
	var res = make([]string, 0)
	for _, t := range tok {
		res = append(res, t)
	}
	return res
}

func tokenError(lit string, expected ...string) error {
	return fmt.Errorf("found %q; expected %s", lit, strings.Join(expected, ", "))
}
