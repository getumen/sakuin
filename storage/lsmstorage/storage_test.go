package lsmstorage_test

import (
	"context"
	"testing"

	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/invertedindex"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/storage/lsmstorage"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {

	num := 100

	ctx := context.Background()

	target, err := lsmstorage.NewStorage(t.TempDir())
	require.NoError(t, err)

	for i := 0; i < num; i++ {
		builder := invertedindex.NewBuilder()
		builder.AddDocument(document.NewDocument(
			uint64(i),
			[]*document.Field{
				document.NewField("title", []term.Term{
					term.NewText("a"),
					term.NewText("b"),
					term.NewText("c"),
					term.NewText("d"),
				}),
				document.NewField("body", []term.Term{
					term.NewText("a"),
					term.NewText("b"),
					term.NewText("c"),
					term.NewText("d"),
				}),
			}))
		index := builder.Build()

		err = target.Merge(ctx, index)
		require.NoError(t, err)
	}

	indexIterator, err := target.GetIndexIterator(ctx, []*termcond.TermCondition{
		termcond.NewEqual(term.NewText("a")),
	})
	require.NoError(t, err)

	require.True(t, indexIterator.HasNext())

	expectedPostings := make([]*posting.Posting, num)
	for i := range expectedPostings {
		expectedPostings[i] = posting.NewPosting(uint64(i), []uint32{0})
	}

	expected := invertedindex.NewInvertedIndex(0)
	expected.Put(term.NewText("a"), fieldindex.NewFieldIndexFromMap(map[fieldname.FieldName]*postinglist.PostingList{
		"title": postinglist.NewPostingList(expectedPostings),
		"body":  postinglist.NewPostingList(expectedPostings),
	}))

	actual := indexIterator.Next()
	require.NoError(t, err)

	require.True(t, expected.Equal(actual))
}
