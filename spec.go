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

type Path interface {
	Segment() string
	SubSegement() Path
}

type pathElement struct {
	segemnt string
	child   Path
}

func (e pathElement) Segment() string {
	return e.segemnt
}

func (e pathElement) SubSegement() Path {
	return e.child
}

type IncludeSpec []Path

type PageSpec struct {
	Offset          int
	Limit           int
	Number          int
	Size            int
	defaultPageSize uint
	maxPageSize     uint
}

func NewPageSpec(defaultPageSize, maxPageSize uint) *PageSpec {
	return &PageSpec{
		Offset:          -1,
		Limit:           -1,
		Number:          -1,
		Size:            -1,
		defaultPageSize: defaultPageSize,
		maxPageSize:     maxPageSize,
	}
}

func (p PageSpec) String() string {
	if p.isOffsetBased() {
		return fmt.Sprintf("%s[%s]=%d", general[PAGE], pageKeywords[OFFSET], p.Offset) +
			general[AMPERSAND] +
			fmt.Sprintf("%s[%s]=%d", general[PAGE], pageKeywords[LIMIT], p.Limit)
	} else if p.isPageBased() {
		return fmt.Sprintf("%s[%s]=%d", general[PAGE], pageKeywords[NUMBER], p.Number) +
			general[AMPERSAND] +
			fmt.Sprintf("%s[%s]=%d", general[PAGE], pageKeywords[SIZE], p.Size)
	} else if p.Offset >= 0 {
		return fmt.Sprintf("%s[%s]=%d", general[PAGE], pageKeywords[OFFSET], p.Offset) +
			general[AMPERSAND] +
			fmt.Sprintf("%s[%s]=%d", general[PAGE], pageKeywords[LIMIT], p.defaultPageSize)
	} else if p.Limit >= 0 {
		return fmt.Sprintf("%s[%s]=%d", general[PAGE], pageKeywords[LIMIT], p.Limit)
	}
	return fmt.Sprintf("%s[%s]=%d", general[PAGE], pageKeywords[LIMIT], p.defaultPageSize)
}

func (p PageSpec) PageOffset() (limit, offset uint) {
	if p.isOffsetBased() {
		offset = uint(p.Offset)
		limit = uint(p.Limit)
	} else if p.isPageBased() {
		offset = uint(p.Size * (p.Number - 1))
		limit = uint(p.Size)
	} else if p.Limit > -1 {
		limit = uint(p.Limit)
	} else {
		limit = p.defaultPageSize
	}
	if limit > p.maxPageSize {
		limit = p.maxPageSize
	}
	return
}

func (p PageSpec) isOffsetBased() bool {
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

func NewQuerySpec(defaultPageSize, maxPageSize uint) *QuerySpec {
	querySpec := new(QuerySpec)
	querySpec.Field = make(FieldSpec)
	querySpec.Filter = make([]*FilterSpec, 0)
	querySpec.PageSpec = NewPageSpec(defaultPageSize, maxPageSize)
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
