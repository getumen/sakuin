package tokenizer_test

import (
	"testing"

	"github.com/getumen/sakuin/analysis/token"
	"github.com/getumen/sakuin/analysis/tokenizer"
	"github.com/getumen/sakuin/term"
	"github.com/stretchr/testify/require"
)

func Test_bigramTokenizer_Tokenize(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want token.TokenStream
	}{
		{
			name: "empty",
			args: args{
				content: "",
			},
			want: token.TokenStream{},
		},
		{
			name: "hello",
			args: args{
				content: "hello",
			},
			want: token.TokenStream{
				token.NewToken(term.NewText("he")),
				token.NewToken(term.NewText("el")),
				token.NewToken(term.NewText("ll")),
				token.NewToken(term.NewText("lo")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tokenizer.NewBigramTokenizer()
			got := tr.Tokenize(tt.args.content)
			require.Equal(t, tt.want, got)
		})
	}
}
