package invertedindex_test

import (
	"testing"

	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/invertedindex"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/term"
	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {
	docs := []*document.Document{
		document.NewDocument(1, []*document.Field{
			document.NewField("title", []term.Term{
				term.NewText("he"),
				term.NewText("el"),
				term.NewText("ll"),
				term.NewText("lo"),
			}),
		}),
		document.NewDocument(2, []*document.Field{
			document.NewField("title", []term.Term{
				term.NewText("he"),
				term.NewText("el"),
				term.NewText("ll"),
			}),
		}),
		document.NewDocument(2, []*document.Field{
			document.NewField("title", []term.Term{
				term.NewText("he"),
				term.NewText("el"),
				term.NewText("ll"),
			}),
		}),
	}

	builder := invertedindex.NewBuilder()

	for _, doc := range docs {
		builder.AddDocument(doc)
	}

	index := builder.Build()

	expected := invertedindex.NewInvertedIndex(0)
	expected.Put(
		term.NewText("he"),
		fieldindex.NewFieldIndexFromMap(map[fieldname.FieldName]*postinglist.PostingList{
			"title": postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, []uint32{0}),
					posting.NewPosting(2, []uint32{0}),
				},
			),
		}),
	)
	expected.Put(
		term.NewText("el"),
		fieldindex.NewFieldIndexFromMap(map[fieldname.FieldName]*postinglist.PostingList{
			"title": postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, []uint32{1}),
					posting.NewPosting(2, []uint32{1}),
				},
			),
		}),
	)
	expected.Put(
		term.NewText("ll"),
		fieldindex.NewFieldIndexFromMap(map[fieldname.FieldName]*postinglist.PostingList{
			"title": postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, []uint32{2}),
					posting.NewPosting(2, []uint32{2}),
				},
			),
		}),
	)
	expected.Put(
		term.NewText("lo"),
		fieldindex.NewFieldIndexFromMap(map[fieldname.FieldName]*postinglist.PostingList{
			"title": postinglist.NewPostingList(
				[]*posting.Posting{
					posting.NewPosting(1, []uint32{3}),
				},
			),
		}),
	)

	expectedIter := expected.Iterator()
	actualIter := index.Iterator()
	for expectedIter.Next() && actualIter.Next() {
		require.Equal(t, expectedIter.Key(), actualIter.Key())
		require.Equal(t, expectedIter.Value(), actualIter.Value())
	}
	require.False(t, expectedIter.Next())
	require.False(t, actualIter.Next())
}
