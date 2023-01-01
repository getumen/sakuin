package expression

import (
	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/termcond"
)

type Expression struct {
	and              []*Expression
	or               []*Expression
	not              *Expression
	phrase           []*Expression
	relativePosition []int64
	feature          *BooleanFeature
}

func NewAnd(arr []*Expression) *Expression {
	return &Expression{
		and: arr,
	}
}

func NewOr(arr []*Expression) *Expression {
	return &Expression{or: arr}
}

func NewNot(value *Expression) *Expression {
	return &Expression{not: value}
}

func NewPhrase(arr []*Expression, relativePosition []int64) *Expression {
	return &Expression{
		phrase:           arr,
		relativePosition: relativePosition,
	}
}

func NewFeature(f *BooleanFeature) *Expression {
	return &Expression{feature: f}
}

func (b Expression) And() []*Expression {
	return b.and
}

func (b Expression) RelativePosition() []int64 {
	return b.relativePosition
}

func (b Expression) Or() []*Expression {
	return b.or
}

func (b Expression) Not() *Expression {
	return b.not
}

func (b Expression) Phrase() []*Expression {
	return b.phrase
}

func (b Expression) Feature() *BooleanFeature {
	return b.feature
}

func (b Expression) TermConditions() []*termcond.TermCondition {
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

	if b.phrase != nil {
		for _, v := range b.phrase {
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
