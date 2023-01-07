package storage_test

import (
	"testing"

	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/storage"
	"github.com/getumen/sakuin/storageutil"
	"github.com/getumen/sakuin/term"
	"github.com/stretchr/testify/require"
)

func TestIndexIterator(t *testing.T) {
	f := fieldindex.NewFieldIndexFromMap(map[fieldname.FieldName]*postinglist.PostingList{
		"f1": postinglist.NewPostingList([]*posting.Posting{
			posting.NewPosting(1, []uint32{1, 2, 4, 5}),
			posting.NewPosting(2, []uint32{1, 2, 4, 5}),
			posting.NewPosting(3, []uint32{1, 2, 4, 5}),
			posting.NewPosting(4, []uint32{1, 2, 4, 5}),
		}),
	})

	seg := storageutil.NewSegment(make([]byte, 0))
	err := seg.Save(1, f)
	require.NoError(t, err)

	target := storage.NewIndexIterator(
		[]term.Term{term.NewText("a"), term.NewText("b"), term.NewText("c")},
		[]*storageutil.SegmentIterator{
			seg.Iterator(),
			seg.Iterator(),
			seg.Iterator(),
		},
	)

	require.True(t, target.HasNext())
	actual, err := target.Next()
	require.NoError(t, err)
	require.False(t, target.HasNext())
	it := actual.Iterator()
	require.True(t, it.Next())
	require.Equal(t, term.NewText("a"), it.Key())
	require.Equal(t, f, it.Value())
	require.True(t, it.Next())
	require.Equal(t, term.NewText("b"), it.Key())
	require.Equal(t, f, it.Value())
	require.True(t, it.Next())
	require.Equal(t, term.NewText("c"), it.Key())
	require.Equal(t, f, it.Value())
	require.False(t, it.Next())
}
