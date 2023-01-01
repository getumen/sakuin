package tokenizer

import (
	"strings"

	"github.com/getumen/sakuin/analysis/token"
	"github.com/getumen/sakuin/term"
)

type bigramTokenizer struct {
}

func NewBigramTokenizer() *bigramTokenizer {
	return &bigramTokenizer{}
}

func (t bigramTokenizer) Tokenize(content string) token.TokenStream {
	result := make(token.TokenStream, 0)

	if len(content) < 2 {
		return result
	}
	chars := strings.Split(content, "")
	for i := 0; i < len(chars)-1; i++ {
		result = append(result, token.NewToken(term.NewText(chars[i]+chars[i+1])))
	}
	return result
}
