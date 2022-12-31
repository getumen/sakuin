package booleanexpression_test

import (
	"testing"

	"github.com/getumen/sakuin/booleanexpression"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
	"github.com/stretchr/testify/require"
)

func TestBooleanExpression_TermConditions(t *testing.T) {
	tests := []struct {
		name string
		b    *booleanexpression.BooleanExpression
		want []*termcond.TermCondition
	}{
		{
			name: "equal",
			b: booleanexpression.NewAnd(
				[]*booleanexpression.BooleanExpression{
					booleanexpression.NewFeature(
						booleanexpression.NewBoolenaFeature(
							"f1",
							termcond.NewEqual(term.NewText("A")),
						),
					),
					booleanexpression.NewFeature(
						booleanexpression.NewBoolenaFeature(
							"f1",
							termcond.NewEqual(term.NewText("B")),
						),
					),
				},
				nil,
			),
			want: []*termcond.TermCondition{
				termcond.NewEqual(term.NewText("A")),
				termcond.NewEqual(term.NewText("B")),
			},
		},
		{
			name: "range",
			b: booleanexpression.NewAnd(
				[]*booleanexpression.BooleanExpression{
					booleanexpression.NewFeature(
						booleanexpression.NewBoolenaFeature(
							"f1",
							termcond.NewRange(
								term.NewInt64(1),
								true,
								term.NewInt64(3),
								false,
							),
						),
					),
					booleanexpression.NewFeature(
						booleanexpression.NewBoolenaFeature(
							"f1",
							termcond.NewRange(
								term.NewInt64(3),
								true,
								term.NewInt64(5),
								false,
							),
						),
					),
				},
				nil,
			),
			want: []*termcond.TermCondition{
				termcond.NewRange(
					term.NewInt64(1),
					true,
					term.NewInt64(5),
					false,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.TermConditions()
			require.Equal(t, tt.want, got)
		})
	}
}
