package termcond

import (
	"sort"

	"github.com/getumen/sakuin/term"
)

type TermCondition struct {
	start        term.Term
	includeStart bool
	end          term.Term
	includeEnd   bool
}

func NewEqual(value term.Term) *TermCondition {
	return &TermCondition{
		start:        value,
		includeStart: true,
		end:          value,
		includeEnd:   true,
	}
}

func (c TermCondition) Start() term.Term {
	return c.start
}

func (c TermCondition) IncludeStart() bool {
	return c.includeStart
}

func (c TermCondition) End() term.Term {
	return c.end
}

func (c TermCondition) IncludeEnd() bool {
	return c.includeEnd
}

func (c TermCondition) IsEqual() bool {
	return term.Comparator(c.start, c.end) == 0 && c.includeStart && c.includeEnd
}

func NewRange(
	start term.Term,
	includeStart bool,
	end term.Term,
	includeEnd bool,
) *TermCondition {
	return &TermCondition{
		start:        start,
		includeStart: includeStart,
		end:          end,
		includeEnd:   includeEnd,
	}
}

func Simplify(conds []*TermCondition) []*TermCondition {
	sort.Slice(conds, func(i, j int) bool {
		if term.Comparator(conds[i].start, conds[j].start) == 0 {
			return conds[i].includeStart
		}
		return term.Comparator(conds[i].start, conds[j].start) < 0
	})

	result := make([]*TermCondition, 0)
	for i, v := range conds {
		if i == 0 {
			result = append(result, v)
			continue
		}
		// no overlap
		if term.Comparator(result[len(result)-1].end, v.start) < 0 ||
			(term.Comparator(result[len(result)-1].end, v.start) == 0 &&
				!result[len(result)-1].includeEnd && !v.includeStart) {
			result = append(result, v)
			continue
		}
		// result[i-1] contains v
		if term.Comparator(v.end, result[len(result)-1].end) < 0 {
			continue
		}
		// ends is equal
		if term.Comparator(v.end, result[len(result)-1].end) == 0 {
			result[len(result)-1].includeEnd = result[len(result)-1].includeEnd || v.includeEnd
			continue
		}
		// partially overlap
		result[len(result)-1].end = v.end
		result[len(result)-1].includeEnd = v.includeEnd
	}
	return result
}
