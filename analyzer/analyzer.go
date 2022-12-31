package analyzer

import "github.com/getumen/sakuin/term"

type CharFilter interface {
	Filter(string) string
}

type Tokenizer interface {
	Tokenize(string) TokenStream
}

type TokenFilter interface {
	Filter(TokenStream) TokenStream
}

type Analyzer struct {
	charFilters  []CharFilter
	tokenizer    Tokenizer
	tokenFilters []TokenFilter
}

type TokenStream struct {
	terms []term.Term
}
