package tokenizer_test

import (
	"testing"

	"github.com/getumen/sakuin/analysis/token"
	"github.com/getumen/sakuin/analysis/tokenizer"
	"github.com/getumen/sakuin/term"
	"github.com/stretchr/testify/require"
)

func Test_japaneseTokenizer(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want token.TokenStream
	}{
		{
			name: "empty string",
			args: args{
				content: "",
			},
			want: make(token.TokenStream, 0),
		},
		{
			name: "すもももももももものうち",
			args: args{
				content: "すもももももももものうち",
			},
			want: token.TokenStream{
				token.NewToken(term.NewText("すもも"), 0, 3),
				token.NewToken(term.NewText("もも"), 4, 6),
				token.NewToken(term.NewText("もも"), 7, 9),
				token.NewToken(term.NewText("うち"), 10, 12),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, err := tokenizer.NewJapaneseTokinizer()
			require.NoError(t, err)

			got := target.Tokenize(tt.args.content)
			require.Equal(t, tt.want, got)
		})
	}
}
