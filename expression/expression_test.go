package expression_test

import (
	"testing"

	"github.com/getumen/sakuin/expression"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
	"github.com/stretchr/testify/require"
)

func TestExpression_TermConditions(t *testing.T) {
	tests := []struct {
		name string
		b    *expression.Expression
		want []*termcond.TermCondition
	}{
		{
			name: "equal",
			b: expression.NewAnd(
				[]*expression.Expression{
					expression.NewFeature(
						expression.NewBoolenaFeature(
							"f1",
							termcond.NewEqual(term.NewText("A")),
						),
					),
					expression.NewFeature(
						expression.NewBoolenaFeature(
							"f1",
							termcond.NewEqual(term.NewText("B")),
						),
					),
				},
			),
			want: []*termcond.TermCondition{
				termcond.NewEqual(term.NewText("A")),
				termcond.NewEqual(term.NewText("B")),
			},
		},
		{
			name: "range",
			b: expression.NewAnd(
				[]*expression.Expression{
					expression.NewFeature(
						expression.NewBoolenaFeature(
							"f1",
							termcond.NewRange(
								term.NewInt64(1),
								true,
								term.NewInt64(3),
								false,
							),
						),
					),
					expression.NewFeature(
						expression.NewBoolenaFeature(
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
