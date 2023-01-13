package lowlevel_test

import (
	"context"
	"testing"

	"github.com/getumen/sakuin/analysis/token"
	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/expression"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/storage/lsmstorage"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
	"github.com/getumen/sakuin/writer"
	"github.com/stretchr/testify/require"
)

func TestSearchTextInLSMTree(t *testing.T) {

	ctx := context.Background()

	storage, err := lsmstorage.NewStorage(t.TempDir())
	require.NoError(t, err)
	defer storage.Close()

	writer := writer.NewIndexWriter(storage)
	require.NoError(t, writer.CreateDocuments(
		ctx,
		[]*document.Document{
			document.NewDocument(1, []*document.Field{
				document.NewField("content", token.TokenStream{
					token.NewToken(term.NewText("i"), 0, 1),
					token.NewToken(term.NewText("am"), 2, 4),
					token.NewToken(term.NewText("a"), 5, 6),
					token.NewToken(term.NewText("pen"), 7, 10),
				}),
			}),
			document.NewDocument(2, []*document.Field{
				document.NewField("content", token.TokenStream{
					token.NewToken(term.NewText("this"), 0, 4),
					token.NewToken(term.NewText("is"), 5, 7),
					token.NewToken(term.NewText("a"), 8, 9),
					token.NewToken(term.NewText("pen"), 10, 13),
				}),
			}),
			document.NewDocument(3, []*document.Field{
				document.NewField("content", token.TokenStream{
					token.NewToken(term.NewText("i"), 0, 1),
					token.NewToken(term.NewText("am"), 2, 4),
					token.NewToken(term.NewText("a"), 5, 6),
					token.NewToken(term.NewText("cat"), 7, 10),
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
				posting.NewPosting(1, []uint32{7}),
				posting.NewPosting(2, []uint32{10}),
			},
		},
		{
			name:  "cat",
			input: "cat",
			expected: []*posting.Posting{
				posting.NewPosting(3, []uint32{7}),
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			query := expression.NewFeature(
				expression.NewFeatureSpec(
					"content", termcond.NewEqual(term.NewText(c.input)),
				))
			indexIterator, err := storage.GetIndexIterator(ctx, query.TermConditions())
			require.NoError(t, err)

			lists := make([]*postinglist.PostingList, 0)

			for indexIterator.HasNext() {
				index := indexIterator.Next()
				require.NoError(t, err)
				lists = append(lists, index.Search(query))
			}

			result := postinglist.Union(lists)

			cursor := result.Cursor()

			for _, c := range c.expected {
				v := cursor.Value()
				require.Equal(t, c, v)
				cursor.Next()
			}
		})
	}
}
