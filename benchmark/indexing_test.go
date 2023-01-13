package benchmark

import (
	"context"
	"math/rand"
	"testing"

	"github.com/getumen/sakuin/analysis/token"
	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/storage/lsmstorage"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/writer"
	"github.com/stretchr/testify/require"
)

func BenchmarkIndexFloat(b *testing.B) {
	length := 1000
	ctx := context.Background()

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
					token.TokenStream{
						token.NewToken(term.NewFloat64(f), 0, 0),
					},
				),
			},
		)
	}

	b.Run("indexing", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			storage, err := lsmstorage.NewStorage(b.TempDir())
			require.NoError(b, err)
			writer := writer.NewIndexWriter(storage)
			require.NoError(b, writer.CreateDocuments(ctx, documents))
		}
	})
}
