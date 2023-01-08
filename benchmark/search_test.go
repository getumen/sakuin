package benchmark

import (
	"context"
	"math/rand"
	"sort"
	"testing"

	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/expression"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/storage/lsmstorage"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
	"github.com/getumen/sakuin/writer"
	"github.com/stretchr/testify/require"
)

func TestSearchRange(t *testing.T) {

	length := 10000
	start := 100
	hit := 10
	end := start + hit
	ctx := context.Background()

	storage, err := lsmstorage.NewStorage(t.TempDir())
	require.NoError(t, err)
	defer storage.Close()

	writer := writer.NewIndexWriter(storage)

	arr := make([]int64, length)
	for i := range arr {
		arr[i] = rand.Int63()
	}

	documents := make([]*document.Document, len(arr))
	for i, f := range arr {
		documents[i] = document.NewDocument(
			uint64(i+1),
			[]*document.Field{
				document.NewField(
					"value",
					[]term.Term{term.NewInt64(f)},
				),
			},
		)
	}

	require.NoError(t, writer.CreateDocuments(ctx, documents))

	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })

	t.Run("search", func(t *testing.T) {

		query := expression.NewFeature(
			expression.NewFeatureSpec(
				"value", termcond.NewRange(
					term.NewInt64(arr[start]),
					true,
					term.NewInt64(arr[end]),
					false,
				),
			))
		it, err := storage.GetIndexIterator(ctx, query.TermConditions())
		require.NoError(t, err)

		lists := make([]*postinglist.PostingList, 0)

		for it.HasNext() {
			value := it.Next()
			if value.Size() == 0 {
				continue
			}
			lists = append(lists, value.Search(query))
		}

		result := postinglist.Union(lists)

		require.Equal(t, hit, result.Len())

	})
}

func BenchmarkSearchRange(b *testing.B) {

	length := 1000000
	start := 100000
	hit := 3000
	end := start + hit
	ctx := context.Background()

	storage, err := lsmstorage.NewStorage(b.TempDir())
	require.NoError(b, err)
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
			it, err := storage.GetIndexIterator(ctx, query.TermConditions())
			require.NoError(b, err)

			lists := make([]*postinglist.PostingList, 0)

			for it.HasNext() {
				value := it.Next()
				if value.Size() == 0 {
					continue
				}
				lists = append(lists, value.Search(query))
			}

			result := postinglist.Union(lists)

			require.Equal(b, hit, result.Len())
		}
	})
}

func TestSearchManyPostings(t *testing.T) {

	rand.Seed(42)

	length := 10
	loop := 10

	ctx := context.Background()

	storage, err := lsmstorage.NewStorage(t.TempDir())
	require.NoError(t, err)
	defer storage.Close()

	writer := writer.NewIndexWriter(storage)

	for x := 0; x < loop; x++ {
		documents := make([]*document.Document, length)
		for i := 0; i < length; i++ {
			if i == 0 {
				documents[i] = document.NewDocument(
					rand.Uint64()%1_000_000,
					[]*document.Field{
						document.NewField(
							"value",
							[]term.Term{
								term.NewText("a"),
								term.NewText("b"),
							},
						),
					},
				)
			} else {
				documents[i] = document.NewDocument(
					rand.Uint64()%1_000_000,
					[]*document.Field{
						document.NewField(
							"value",
							[]term.Term{
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
							},
						),
					},
				)
			}

		}
		require.NoError(t, writer.CreateDocuments(ctx, documents))
	}

	t.Run("search", func(t *testing.T) {
		query := expression.NewPhrase(
			[]*expression.Expression{
				expression.NewFeature(
					expression.NewFeatureSpec(
						"value", termcond.NewEqual(term.NewText("a")),
					)),
				expression.NewFeature(
					expression.NewFeatureSpec(
						"value", termcond.NewEqual(term.NewText("b")),
					)),
			},
			[]uint32{0, 1},
		)
		it, err := storage.GetIndexIterator(ctx, query.TermConditions())
		require.NoError(t, err)

		lists := make([]*postinglist.PostingList, 0)

		for it.HasNext() {
			value := it.Next()
			if value.Size() == 0 {
				continue
			}
			lists = append(lists, value.Search(query))
		}

		result := postinglist.Union(lists)

		require.Equal(t, loop, result.Len())
	})
}

func BenchmarkSearchManyPostings(b *testing.B) {

	rand.Seed(42)

	length := 10000
	loop := 10

	ctx := context.Background()

	storage, err := lsmstorage.NewStorage(b.TempDir())
	require.NoError(b, err)
	defer storage.Close()

	writer := writer.NewIndexWriter(storage)

	for x := 0; x < loop; x++ {
		documents := make([]*document.Document, length)
		for i := 0; i < length; i++ {
			if i == 0 {
				documents[i] = document.NewDocument(
					rand.Uint64()%1_000_000_000,
					[]*document.Field{
						document.NewField(
							"value",
							[]term.Term{
								term.NewText("a"),
								term.NewText("b"),
							},
						),
					},
				)
			} else {
				documents[i] = document.NewDocument(
					rand.Uint64()%1_000_000_000,
					[]*document.Field{
						document.NewField(
							"value",
							[]term.Term{
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
								term.NewText("a"),
							},
						),
					},
				)
			}

		}
		require.NoError(b, writer.CreateDocuments(ctx, documents))
	}

	b.Run("search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			query := expression.NewPhrase(
				[]*expression.Expression{
					expression.NewFeature(
						expression.NewFeatureSpec(
							"value", termcond.NewEqual(term.NewText("a")),
						)),
					expression.NewFeature(
						expression.NewFeatureSpec(
							"value", termcond.NewEqual(term.NewText("b")),
						)),
				},
				[]uint32{0, 1},
			)
			it, err := storage.GetIndexIterator(ctx, query.TermConditions())
			require.NoError(b, err)

			lists := make([]*postinglist.PostingList, 0)

			for it.HasNext() {
				value := it.Next()
				if value.Size() == 0 {
					continue
				}
				lists = append(lists, value.Search(query))
			}

			result := postinglist.Union(lists)

			require.Equal(b, loop, result.Len())
		}
	})
}
