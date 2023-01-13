package token

import (
	"fmt"

	"github.com/getumen/sakuin/term"
)

type Token struct {
	term term.Term
	// start specifies the byte offset of the beginning of the term in the field.
	start uint32
	// end specifies the byte offset of the end of the term in the field.
	end uint32
}

func NewToken(
	term term.Term,
	start uint32,
	end uint32,
) *Token {
	return &Token{
		term:  term,
		start: start,
		end:   end,
	}
}

func (t *Token) Term() term.Term {
	return t.term
}

func (t *Token) Start() uint32 {
	return t.start
}

func (t *Token) End() uint32 {
	return t.end
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"token: %s  Type: %s",
		t.term.String(),
		t.term.Type(),
	)
}

type TokenStream []*Token

func (t TokenStream) Terms() []term.Term {
	result := make([]term.Term, len(t))
	for i := range t {
		result[i] = t[i].term
	}
	return result
}
