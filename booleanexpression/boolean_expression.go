package booleanexpression

import (
	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/termcond"
)

type BooleanExpression struct {
	and              []*BooleanExpression
	or               []*BooleanExpression
	not              *BooleanExpression
	relativePosition []int64
	feature          *BooleanFeature
}

func NewAnd(arr []*BooleanExpression, relativePosition []int64) *BooleanExpression {
	return &BooleanExpression{and: arr, relativePosition: relativePosition}
}

func NewOr(arr []*BooleanExpression) *BooleanExpression {
	return &BooleanExpression{or: arr}
}

func NewNot(value *BooleanExpression) *BooleanExpression {
	return &BooleanExpression{not: value}
}

func NewFeature(f *BooleanFeature) *BooleanExpression {
	return &BooleanExpression{feature: f}
}

func (b BooleanExpression) And() []*BooleanExpression {
	return b.and
}

func (b BooleanExpression) RelativePosition() []int64 {
	return b.relativePosition
}

func (b BooleanExpression) Or() []*BooleanExpression {
	return b.or
}

func (b BooleanExpression) Not() *BooleanExpression {
	return b.not
}

func (b BooleanExpression) Feature() *BooleanFeature {
	return b.feature
}

func (b BooleanExpression) TermConditions() []*termcond.TermCondition {
	if b.feature != nil {
		return []*termcond.TermCondition{b.feature.termCondition}
	}
	if b.not != nil {
		return b.not.TermConditions()
	}

	conds := make([]*termcond.TermCondition, 0)
	if b.and != nil {
		for _, v := range b.and {
			conds = append(conds, v.TermConditions()...)
		}
	}

	if b.or != nil {
		for _, v := range b.or {
			conds = append(conds, v.TermConditions()...)
		}
	}

	return termcond.Simplify(conds)
}

type BooleanFeature struct {
	field         fieldname.FieldName
	termCondition *termcond.TermCondition
}

func NewBoolenaFeature(
	field fieldname.FieldName,
	termCondition *termcond.TermCondition,
) *BooleanFeature {
	return &BooleanFeature{
		field:         field,
		termCondition: termCondition,
	}
}

func (f BooleanFeature) TermCondition() *termcond.TermCondition {
	return f.termCondition
}

func (f BooleanFeature) Field() fieldname.FieldName {
	return f.field
}
