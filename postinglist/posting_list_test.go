package postinglist_test

import (
	"testing"

	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/stretchr/testify/require"
)

func TestPhraseMatch(t *testing.T) {
	type args struct {
		postingLists     []*postinglist.PostingList
		relativePosition []uint32
	}
	tests := []struct {
		name string
		args args
		want *postinglist.PostingList
	}{
		{
			name: "phrase match 1",
			args: args{
				postingLists: []*postinglist.PostingList{
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, []uint32{1}),
					}),
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, []uint32{3, 7}),
					}),
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, []uint32{5}),
					}),
				},
				relativePosition: []uint32{0, 2, 4},
			},
			want: postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, []uint32{1}),
				},
			),
		},
		{
			name: "phrase match 2",
			args: args{
				postingLists: []*postinglist.PostingList{
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, []uint32{1, 5, 7}),
						posting.NewPosting(2, []uint32{3, 5, 7}),
					}),
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, []uint32{2}),
						posting.NewPosting(2, []uint32{2}),
					}),
				},
				relativePosition: []uint32{0, 1},
			},
			want: postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, []uint32{1}),
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := postinglist.PhraseMatch(
				tt.args.postingLists,
				tt.args.relativePosition,
			)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIntersection(t *testing.T) {
	type args struct {
		postingLists []*postinglist.PostingList
	}
	tests := []struct {
		name string
		args args
		want *postinglist.PostingList
	}{
		{
			name: "intersection",
			args: args{
				postingLists: []*postinglist.PostingList{
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, []uint32{1}),
					}),
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, []uint32{3, 7}),
					}),
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, []uint32{5}),
					}),
				},
			},
			want: postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, []uint32{1, 3, 5, 7}),
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := postinglist.Intersection(tt.args.postingLists)
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
						posting.NewPosting(1, []uint32{1}),
						posting.NewPosting(2, []uint32{1}),
						posting.NewPosting(4, []uint32{1, 3, 4}),
					}),
					postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(3, []uint32{1}),
						posting.NewPosting(4, []uint32{2, 4}),
					}),
				},
			},
			want: postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, []uint32{1}),
					posting.NewPosting(2, []uint32{1}),
					posting.NewPosting(3, []uint32{1}),
					posting.NewPosting(4, []uint32{1, 2, 3, 4}),
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
			name: "{1,3,5} - {2,3}",
			args: args{
				x: postinglist.NewPostingList([]*posting.Posting{
					posting.NewPosting(1, []uint32{1}),
					posting.NewPosting(3, []uint32{1}),
					posting.NewPosting(5, []uint32{1}),
				}),
				y: postinglist.NewPostingList([]*posting.Posting{
					posting.NewPosting(2, []uint32{1}),
					posting.NewPosting(3, []uint32{1}),
				}),
			},
			want: postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, []uint32{1}),
					posting.NewPosting(5, []uint32{1}),
				},
			),
		},
		{
			name: "{1} - {}",
			args: args{
				x: postinglist.NewPostingList([]*posting.Posting{
					posting.NewPosting(1, []uint32{1}),
				}),
				y: postinglist.NewPostingList([]*posting.Posting{}),
			},
			want: postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, []uint32{1}),
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
