package postinglist_test

import (
	"testing"

	"github.com/getumen/sakuin/position"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/stretchr/testify/require"
)

func TestIntersection(t *testing.T) {
	type args struct {
		postingLists     []*postinglist.PostingList
		relativePosition []int64
	}
	tests := []struct {
		name string
		args args
		want *postinglist.PostingList
	}{
		{
			name: "intersection relativePositionなし",
			args: args{
				postingLists: []*postinglist.PostingList{
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, nil),
					}),
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, nil),
					}),
				},
				relativePosition: nil,
			},
			want: postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, nil),
				},
			),
		},
		{
			name: "intersection relativePositionあり positionもマッチ",
			args: args{
				postingLists: []*postinglist.PostingList{
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, position.NewPositions([]int64{1, 8, 10})),
					}),
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, position.NewPositions([]int64{3, 10})),
					}),
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, position.NewPositions([]int64{5, 12})),
					}),
				},
				relativePosition: []int64{0, 2, 4},
			},
			want: postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, nil),
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := postinglist.Intersection(tt.args.postingLists, tt.args.relativePosition)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnion(t *testing.T) {
	type args struct {
		postingLists []*postinglist.PostingList
	}
	tests := []struct {
		name string
		args args
		want *postinglist.PostingList
	}{
		{
			name: "union",
			args: args{
				postingLists: []*postinglist.PostingList{
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, nil),
						posting.NewPosting(2, nil),
						posting.NewPosting(4, nil),
					}),
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(3, nil),
						posting.NewPosting(4, nil),
					}),
				},
			},
			want: postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, nil),
					posting.NewPosting(2, nil),
					posting.NewPosting(3, nil),
					posting.NewPosting(4, nil),
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := postinglist.Union(tt.args.postingLists)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDifference(t *testing.T) {
	type args struct {
		x *postinglist.PostingList
		y *postinglist.PostingList
	}
	tests := []struct {
		name string
		args args
		want *postinglist.PostingList
	}{
		{
			name: "union",
			args: args{
				x: postinglist.NewPostingList([]*posting.Posting{
					posting.NewPosting(1, nil),
					posting.NewPosting(3, nil),
				}),
				y: postinglist.NewPostingList([]*posting.Posting{
					posting.NewPosting(2, nil),
					posting.NewPosting(3, nil),
				}),
			},
			want: postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, nil),
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := postinglist.Difference(tt.args.x, tt.args.y)
			require.Equal(t, tt.want, got)
		})
	}
}
