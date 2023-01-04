package benchmark

import (
	"context"
	"math/rand"
	"sort"
	"testing"

	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/expression"
	"github.com/getumen/sakuin/storage/memstorage"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
	"github.com/getumen/sakuin/writer"
	"github.com/stretchr/testify/require"
)

func BenchmarkSearchRange(b *testing.B) {
	length := 1000000
	start := 100000
	hit := 10
	end := start + hit
	ctx := context.Background()

	storage := memstorage.NewMemStorage()
	defer storage.Close()

	writer := writer.NewIndexWriter(storage)

	arr := make([]float64, length)
	for i := range arr {
		arr[i] = rand.Float64()
	}

	documents := make([]*document.Document, len(arr))
	for i, f := range arr {
		documents[i] = document.NewDocument(
			uint64(i+1),
			[]*document.Field{
				document.NewField(
					"value",
					[]term.Term{term.NewFloat64(f)},
				),
			},
		)
	}

	require.NoError(b, writer.CreateDocuments(ctx, documents))

	sort.Float64s(arr)

	b.Run("search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			query := expression.NewFeature(
				expression.NewFeatureSpec(
					"value", termcond.NewRange(
						term.NewFloat64(arr[start]),
						true,
						term.NewFloat64(arr[end]),
						false,
					),
				))
			partialIndex, err := storage.GetIndex(ctx, query.TermConditions())
			require.NoError(b, err)

			result := partialIndex.Search(query)
			require.Equal(b, hit, result.Len())
		}
	})

}
