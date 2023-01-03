package posting_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/getumen/sakuin/posting"
	"github.com/stretchr/testify/require"
)

func TestSerializeDeserialize(t *testing.T) {

	const docNum = 10
	const positionNum = 10

	postingList := make([]*posting.Posting, docNum)

	docIDs := make([]uint64, docNum)
	for i := 0; i < docNum; i++ {
		docIDs[i] = rand.Uint64() % (1 << 50)
	}
	sort.Slice(docIDs, func(i, j int) bool {
		return docIDs[i] < docIDs[j]
	})

	for i, v := range docIDs {
		positionList := make([]uint32, positionNum)
		for p := 0; p < positionNum; p++ {
			positionList[p] = rand.Uint32()
		}
		sort.Slice(positionList, func(i, j int) bool { return positionList[i] < positionList[j] })
		postingList[i] = posting.NewPosting(v, positionList)
	}

	b, err := posting.Serialize(postingList)
	require.NoError(t, err)
	actual, err := posting.Deserialize(b)
	require.NoError(t, err)

	require.Equal(t, postingList, actual)
}
