package analysis

import (
	"github.com/getumen/sakuin/analysis/token"
)

type CharFilter interface {
	Filter(string) string
}

type Tokenizer interface {
	Tokenize(string) token.TokenStream
}

type TokenFilter interface {
	Filter(token.TokenStream) token.TokenStream
}

type Analyzer struct {
	charFilters  []CharFilter
	tokenizer    Tokenizer
	tokenFilters []TokenFilter
}
