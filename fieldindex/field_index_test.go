package fieldindex_test

import (
	"testing"

	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/stretchr/testify/require"
)

func TestFieldIndex_Merge(t *testing.T) {
	type args struct {
		other fieldindex.FieldIndex
	}
	tests := []struct {
		name     string
		f        fieldindex.FieldIndex
		args     args
		expected fieldindex.FieldIndex
	}{
		{
			name: "異なるフィールドのマージ",
			f: map[fieldname.FieldName]*postinglist.PostingList{
				"a": postinglist.NewPostingList([]*posting.Posting{
					posting.NewPosting(1, nil),
				}),
			},
			args: args{
				other: map[fieldname.FieldName]*postinglist.PostingList{
					"b": postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(1, nil),
					}),
				},
			},
			expected: map[fieldname.FieldName]*postinglist.PostingList{
				"a": postinglist.NewPostingList([]*posting.Posting{
					posting.NewPosting(1, nil),
				}),
				"b": postinglist.NewPostingList([]*posting.Posting{
					posting.NewPosting(1, nil),
				}),
			},
		},
		{
			name: "同じフィールドのマージ",
			f: map[fieldname.FieldName]*postinglist.PostingList{
				"a": postinglist.NewPostingList([]*posting.Posting{
					posting.NewPosting(1, nil),
				}),
			},
			args: args{
				other: map[fieldname.FieldName]*postinglist.PostingList{
					"a": postinglist.NewPostingList([]*posting.Posting{
						posting.NewPosting(2, nil),
					}),
				},
			},
			expected: map[fieldname.FieldName]*postinglist.PostingList{
				"a": postinglist.NewPostingList([]*posting.Posting{
					posting.NewPosting(1, nil),
					posting.NewPosting(2, nil),
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.Merge(tt.args.other)
			require.Equal(t, tt.expected, tt.f)
		})
	}
}
