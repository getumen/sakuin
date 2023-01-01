package benchmark

import (
	"context"
	"math/rand"
	"testing"

	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/storage/memstorage"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/writer"
	"github.com/stretchr/testify/require"
)

func BenchmarkIndexFloat(b *testing.B) {
	length := 1000000
	ctx := context.Background()

	arr := make([]float64, length)
	for i := range arr {
		arr[i] = rand.Float64()
	}

	documents := make([]*document.Document, len(arr))
	for i, f := range arr {
		documents[i] = document.NewDocument(
			int64(i+1),
			[]*document.Field{
				document.NewField(
					"value",
					[]term.Term{term.NewFloat64(f)},
				),
			},
		)
	}

	b.Run("indexing", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			storage := memstorage.NewMemStorage()
			writer := writer.NewIndexWriter(storage)
			require.NoError(b, writer.CreateDocuments(ctx, documents))
		}
	})

}
