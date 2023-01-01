package tokenizer

import (
	"fmt"

	"github.com/getumen/sakuin/analysis/token"
	"github.com/getumen/sakuin/term"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type japaneseTokenizer struct {
	t         *tokenizer.Tokenizer
	posFilter *filter.POSFilter
}

func NewJapaneseTokinizer() (*japaneseTokenizer, error) {
	t, err := tokenizer.New(ipa.DictShrink(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, fmt.Errorf("fail to create tokenizer: %w", err)
	}
	posFilter := filter.NewPOSFilter(filter.POS{"助詞"}, filter.POS{"記号"})
	return &japaneseTokenizer{
		t:         t,
		posFilter: posFilter,
	}, err
}

func (t japaneseTokenizer) Tokenize(content string) token.TokenStream {
	tokens := t.t.Analyze(content, tokenizer.Search)
	t.posFilter.Drop(&tokens)
	result := make(token.TokenStream, 0)
	for _, s := range tokens {
		result = append(result, token.NewToken(term.NewText(s.Surface)))
	}
	return result
}
