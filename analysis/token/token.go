package token

import (
	"fmt"

	"github.com/getumen/sakuin/term"
)

type Token struct {
	term term.Term
}

func NewToken(term term.Term) *Token {
	return &Token{
		term: term,
	}
}

func (t *Token) Term() term.Term {
	return t.term
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"oken: %s  Type: %s",
		t.term.String(),
		t.term.Type(),
	)
}

type TokenStream []*Token