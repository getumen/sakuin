package storageutil_test

import (
	"math/rand"
	"testing"

	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/storageutil"
	"github.com/stretchr/testify/require"
)

func TestSegment(t *testing.T) {

	segmentNum := 3

	target := storageutil.NewSegment(make([]byte, 0))

	expectd := map[uint64]struct{}{}

	for i := uint64(0); i < 10; i++ {

		expectd[i] = struct{}{}

		index := fieldindex.NewFieldIndexFromMap(
			map[fieldname.FieldName]*postinglist.PostingList{
				"f1": postinglist.NewPostingList([]*posting.Posting{
					posting.NewPosting(i, []uint32{1, 2, 3}),
				}),
			},
		)

		segmentID := int(rand.Uint32()%uint32(segmentNum)) + 1

		db, err := target.Get(segmentID)
		require.NoError(t, err)
		db.Merge(index)

		err = target.Save(segmentID, db)
		require.NoError(t, err)
	}

	actual := map[uint64]struct{}{}

	for i := 1; i <= segmentNum; i++ {
		index, err := target.Get(i)
		require.NoError(t, err)
		for _, value := range index {
			cur := value.Cursor()
			for {
				pos := cur.Value()
				actual[pos.DocIDForTest()] = struct{}{}
				if !cur.Next() {
					break
				}
			}
		}
	}

	require.Equal(t, expectd, actual)

	actual = map[uint64]struct{}{}

	it := target.Iterator()
	for it.HasNext() {
		index, err := it.Next()
		require.NoError(t, err)
		for _, value := range index {
			cur := value.Cursor()
			for {
				pos := cur.Value()
				actual[pos.DocIDForTest()] = struct{}{}
				if !cur.Next() {
					break
				}
			}
		}
	}

	require.Equal(t, expectd, actual)
}
