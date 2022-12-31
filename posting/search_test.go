package posting_test

import (
	"testing"

	"github.com/getumen/sakuin/posting"
	"github.com/stretchr/testify/require"
)

func TestBinarySearch(t *testing.T) {
	type args struct {
		sortedArray []*posting.Posting
		docID       int64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "found",
			args: args{
				sortedArray: []*posting.Posting{
					posting.NewPosting(1, nil),
					posting.NewPosting(3, nil),
					posting.NewPosting(6, nil),
					posting.NewPosting(10, nil),
					posting.NewPosting(15, nil),
					posting.NewPosting(21, nil),
					posting.NewPosting(28, nil),
					posting.NewPosting(36, nil),
					posting.NewPosting(45, nil),
					posting.NewPosting(55, nil),
				},
				docID: 6,
			},
			want: 2,
		},
		{
			name: "not found",
			args: args{
				sortedArray: []*posting.Posting{
					posting.NewPosting(1, nil),
					posting.NewPosting(3, nil),
					posting.NewPosting(6, nil),
					posting.NewPosting(10, nil),
					posting.NewPosting(15, nil),
					posting.NewPosting(21, nil),
					posting.NewPosting(28, nil),
					posting.NewPosting(36, nil),
					posting.NewPosting(45, nil),
					posting.NewPosting(55, nil),
				},
				docID: 14,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := posting.BinarySearch(tt.args.sortedArray, tt.args.docID)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestExponentialSearch(t *testing.T) {
	type args struct {
		sortedArray []*posting.Posting
		docID       int64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "found",
			args: args{
				sortedArray: []*posting.Posting{
					posting.NewPosting(1, nil),
					posting.NewPosting(3, nil),
					posting.NewPosting(6, nil),
					posting.NewPosting(10, nil),
					posting.NewPosting(15, nil),
					posting.NewPosting(21, nil),
					posting.NewPosting(28, nil),
					posting.NewPosting(36, nil),
					posting.NewPosting(45, nil),
					posting.NewPosting(55, nil),
				},
				docID: 6,
			},
			want: 2,
		},
		{
			name: "not found",
			args: args{
				sortedArray: []*posting.Posting{
					posting.NewPosting(1, nil),
					posting.NewPosting(3, nil),
					posting.NewPosting(6, nil),
					posting.NewPosting(10, nil),
					posting.NewPosting(15, nil),
					posting.NewPosting(21, nil),
					posting.NewPosting(28, nil),
					posting.NewPosting(36, nil),
					posting.NewPosting(45, nil),
					posting.NewPosting(55, nil),
				},
				docID: 14,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := posting.ExponentialSearch(tt.args.sortedArray, tt.args.docID)
			require.Equal(t, tt.want, got)
		})
	}
}
