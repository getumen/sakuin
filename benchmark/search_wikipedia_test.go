package benchmark_test

import (
	"context"
	"path"
	"testing"

	"github.com/getumen/sakuin/analysis"
	"github.com/getumen/sakuin/analysis/charfilter"
	"github.com/getumen/sakuin/analysis/tokenizer"
	"github.com/getumen/sakuin/expression"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/storage/lsmstorage"
	"github.com/getumen/sakuin/termcond"
	"github.com/stretchr/testify/require"
)

const tmpDir = "/tmp/sakuin/"

func BenchmarkSearchWikipedia(b *testing.B) {

	indexPath := path.Join(tmpDir, "index")

	analyzer := analysis.NewAnalyzer(
		[]analysis.CharFilter{
			charfilter.NewUnicodeNFKCFilter(),
		},
		tokenizer.NewBigramTokenizer(),
		[]analysis.TokenFilter{},
	)

	storage, err := lsmstorage.NewStorage(indexPath)
	require.NoError(b, err)
	defer storage.Close()

	var hits, segmentNum int
	docs := make(map[uint64]struct{})

	b.Run("search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx := context.Background()

			tokenStream := analyzer.Analyze("日本")
			terms := tokenStream.Terms()
			titlePhrase := make([]*expression.Expression, len(terms))
			bodyPhrase := make([]*expression.Expression, len(terms))
			relativePosition := make([]uint32, len(terms))
			for i := range terms {
				titlePhrase[i] = expression.NewFeature(expression.NewFeatureSpec(
					"title", termcond.NewEqual(terms[i])))
				bodyPhrase[i] = expression.NewFeature(expression.NewFeatureSpec(
					"body", termcond.NewEqual(terms[i])))
				relativePosition[i] = uint32(i)
			}

			query := expression.NewOr(
				[]*expression.Expression{
					expression.NewPhrase(titlePhrase, relativePosition),
					expression.NewPhrase(titlePhrase, relativePosition),
				},
			)
			it, err := storage.GetIndexIterator(ctx, query.TermConditions())
			require.NoError(b, err)

			lists := make([]*postinglist.PostingList, 0)

			var segNum int
			for it.HasNext() {
				value := it.Next()
				lists = append(lists, value.Search(query))
				segNum++
			}

			result := postinglist.Union(lists)
			hits = result.Len()
			segmentNum = segNum
			cur := result.Cursor()
			for {
				docs[cur.Value().DocIDForTest()] = struct{}{}
				if !cur.Next() {
					break
				}
			}
		}
	})

	b.Logf("%d hits in %d segment\n", hits, segmentNum)
	for docID := range docs {
		b.Logf("https://ja.wikipedia.org/w/index.php?curid=%d", docID)
	}
}
