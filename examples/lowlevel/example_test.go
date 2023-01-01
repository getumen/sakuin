package lowlevel_test

import (
	"context"
	"testing"

	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/expression"
	"github.com/getumen/sakuin/position"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/storage/memstorage"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
	"github.com/getumen/sakuin/writer"
	"github.com/stretchr/testify/require"
)

func TestSearchText(t *testing.T) {

	ctx := context.Background()

	storage := memstorage.NewMemStorage()

	writer := writer.NewIndexWriter(storage)
	require.NoError(t, writer.CreateDocuments(
		ctx,
		[]*document.Document{
			document.NewDocument(1, []*document.Field{
				document.NewField("content", []term.Term{
					term.NewText("i"),
					term.NewText("am"),
					term.NewText("a"),
					term.NewText("pen"),
				}),
			}),
			document.NewDocument(2, []*document.Field{
				document.NewField("content", []term.Term{
					term.NewText("this"),
					term.NewText("is"),
					term.NewText("a"),
					term.NewText("pen"),
				}),
			}),
			document.NewDocument(3, []*document.Field{
				document.NewField("content", []term.Term{
					term.NewText("i"),
					term.NewText("am"),
					term.NewText("a"),
					term.NewText("cat"),
				}),
			}),
		}),
	)

	for _, c := range []struct {
		name     string
		input    string
		expected []*posting.Posting
	}{
		{
			name:  "pen",
			input: "pen",
			expected: []*posting.Posting{
				posting.NewPosting(1, position.NewPositions([]int64{3})),
				posting.NewPosting(2, position.NewPositions([]int64{3})),
			},
		},
		{
			name:  "cat",
			input: "cat",
			expected: []*posting.Posting{
				posting.NewPosting(3, position.NewPositions([]int64{3})),
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			query := expression.NewFeature(expression.NewFeatureSpec(
				"content", termcond.NewEqual(term.NewText(c.input)),
			))
			partialIndex, err := storage.GetIndex(ctx, query.TermConditions())
			require.NoError(t, err)

			result := partialIndex.Search(query)
			cursor := result.Cursor()

			for _, c := range c.expected {
				v := cursor.Value()
				require.Equal(t, c, v)
				cursor.Next()
			}
		})
	}

}
