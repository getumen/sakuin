package termcond

import (
	"reflect"
	"testing"

	"github.com/getumen/sakuin/term"
)

func TestSimplify(t *testing.T) {
	type args struct {
		conds []*TermCondition
	}
	tests := []struct {
		name string
		args args
		want []*TermCondition
	}{
		{
			name: "int continuous",
			args: args{
				conds: []*TermCondition{
					NewRange(
						term.NewInt64(1),
						true,
						term.NewInt64(3),
						false,
					),
					NewRange(
						term.NewInt64(3),
						true,
						term.NewInt64(5),
						false,
					),
				},
			},
			want: []*TermCondition{
				NewRange(
					term.NewInt64(1),
					true,
					term.NewInt64(5),
					false,
				),
			},
		},
		{
			name: "float contain",
			args: args{
				conds: []*TermCondition{
					NewRange(
						term.NewInt64(3),
						true,
						term.NewInt64(5),
						false,
					),
					NewRange(
						term.NewInt64(1),
						true,
						term.NewInt64(10),
						false,
					),
				},
			},
			want: []*TermCondition{
				NewRange(
					term.NewInt64(1),
					true,
					term.NewInt64(10),
					false,
				),
			},
		},
		{
			name: "float same",
			args: args{
				conds: []*TermCondition{
					NewRange(
						term.NewInt64(3),
						true,
						term.NewInt64(5),
						false,
					),
					NewRange(
						term.NewInt64(3),
						true,
						term.NewInt64(5),
						false,
					),
				},
			},
			want: []*TermCondition{
				NewRange(
					term.NewInt64(3),
					true,
					term.NewInt64(5),
					false,
				),
			},
		},
		{
			name: "float partial overlap",
			args: args{
				conds: []*TermCondition{
					NewRange(
						term.NewInt64(1),
						true,
						term.NewInt64(4),
						false,
					),
					NewRange(
						term.NewInt64(3),
						true,
						term.NewInt64(5),
						false,
					),
				},
			},
			want: []*TermCondition{
				NewRange(
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
			if got := Simplify(tt.args.conds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Simplify() = %v, want %v", got, tt.want)
			}
		})
	}
}
