package benchmark_test

import (
	"context"
	"path"
	"testing"

	"github.com/getumen/sakuin/analysis"
	"github.com/getumen/sakuin/analysis/charfilter"
	"github.com/getumen/sakuin/analysis/tokenizer"
	"github.com/getumen/sakuin/expression"
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

	storage, err := lsmstorage.NewLSMStorage(indexPath)
	require.NoError(b, err)
	defer storage.Close()

	var hits int

	for i := 0; i < b.N; i++ {
		ctx := context.Background()

		tokenStream := analyzer.Analyze("奈良県")
		terms := tokenStream.Terms()
		phrase := make([]*expression.Expression, len(terms))
		relativePosition := make([]uint32, len(terms))
		for i := range terms {
			phrase[i] = expression.NewFeature(expression.NewFeatureSpec(
				"title", termcond.NewEqual(terms[i])))
			relativePosition[i] = uint32(i)
		}

		query := expression.NewPhrase(phrase, relativePosition)
		partialIndex, err := storage.GetIndex(ctx, query.TermConditions())
		require.NoError(b, err)

		result := partialIndex.Search(query)
		hits = result.Len()
	}
	b.Logf("%d hits\n", hits)
}
