package querystring

import (
	"fmt"
	"strings"
)

type FieldSpec map[string][]string

func (f FieldSpec) String() string {
	if len(f) == 0 {
		return ""
	}
	specs := make([]string, 0)
	for k, v := range f {
		specs = append(
			specs,
			fmt.Sprintf("%s[%s]=%s", keywords[FIELD], k, strings.Join(v, general[COMMA])),
		)
	}
	return strings.Join(specs, general[AMPERSAND])
}

type FilterSpec struct {
	Field    string
	Operator string
	Value    string
}

func (f FilterSpec) String() string {
	return fmt.Sprintf("%s[%s][%s]=%s", general[FILTER], f.Field, f.Operator, f.Value)
}

type IncludeSpec []string

type PageSpec struct {
	Offset int
	Limit  int
	Number int
	Size   int
}

func NewPageSpec() *PageSpec {
	return &PageSpec{Offset: -1, Limit: -1, Number: -1, Size: -1}
}

func (p PageSpec) String() string {
	if p.isPointerBased() {
		return fmt.Sprintf("%s[%s]=%d", general[PAGE], general[OFFSET], p.Offset) +
			general[AMPERSAND] +
			fmt.Sprintf("%s[%s]=%d", general[PAGE], general[LIMIT], p.Limit)
	} else if p.isPageBased() {
		return fmt.Sprintf("%s[%s]=%d", general[PAGE], general[NUMBER], p.Number) +
			general[AMPERSAND] +
			fmt.Sprintf("%s[%s]=%d", general[PAGE], general[SIZE], p.Size)
	}
	return ""
}

func (p PageSpec) Pointer() (limit, offset int) {
	if p.isPointerBased() {
		offset = p.Offset
		limit = p.Limit
	} else if p.isPageBased() {
		offset = p.Size * (p.Number - 1)
		limit = p.Size
	}
	return
}

func (p PageSpec) isPointerBased() bool {
	return p.Offset > -1 && p.Limit > -1
}

func (p PageSpec) isPageBased() bool {
	return p.Number > -1 && p.Size > -1
}

type SortSpec struct {
	Field string
	Desc  bool
}

func (s SortSpec) String() string {
	desc := ""
	if s.Desc {
		desc = general[DASH]
	}
	return fmt.Sprintf("%s%s", desc, s.Field)
}

type QuerySpec struct {
	Field       FieldSpec
	Filter      []*FilterSpec
	IncludeSpec IncludeSpec
	PageSpec    *PageSpec
	Sort        []*SortSpec
}

func NewQuerySpec() *QuerySpec {
	querySpec := new(QuerySpec)
	querySpec.Field = make(FieldSpec)
	querySpec.Filter = make([]*FilterSpec, 0)
	querySpec.PageSpec = NewPageSpec()
	querySpec.Sort = make([]*SortSpec, 0)
	return querySpec
}

func (q QuerySpec) String() string {
	specs := make([]string, 0)
	if field := q.Field.String(); field != "" {
		specs = append(specs, field)
	}

	for _, filter := range q.Filter {
		specs = append(specs, filter.String())
	}

	if page := q.PageSpec.String(); page != "" {
		specs = append(specs, page)
	}

	sortSpecs := make([]string, 0)
	for _, sort := range q.Sort {
		sortSpecs = append(sortSpecs, sort.String())
	}
	if len(sortSpecs) > 0 {
		specs = append(specs, fmt.Sprintf("%s=%s", general[SORT], strings.Join(sortSpecs, ",")))
	}

	return strings.Join(specs, general[AMPERSAND])
}
