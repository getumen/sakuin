package tokenfilter

import (
	"github.com/getumen/sakuin/analysis/token"
	"github.com/getumen/sakuin/term"
)

type japaneseStopwords struct {
	stopWords map[string]struct{}
}

func NewJapanaseStopwords() *japaneseStopwords {
	// https://github.com/stopwords-iso/stopwords-ja
	stopWordsSrc := []string{
		"あそこ", "あっ", "あの", "あのかた", "あの人", "あり", "あります", "ある", "あれ", "い", "いう", "います", "いる", "う", "うち",
		"え", "お", "および", "おり", "おります", "か", "かつて", "から", "が", "き", "ここ", "こちら", "こと", "この", "これ", "これら",
		"さ", "さらに", "し", "しかし", "する", "ず", "せ", "せる", "そこ", "そして", "その", "その他", "その後", "それ", "それぞれ",
		"それで", "た", "ただし", "たち", "ため", "たり", "だ", "だっ", "だれ", "つ", "て", "で", "でき", "できる", "です", "では", "でも",
		"と", "という", "といった", "とき", "ところ", "として", "とともに", "とも", "と共に", "どこ", "どの", "な", "ない", "なお",
		"なかっ", "ながら", "なく", "なっ", "など", "なに", "なら", "なり", "なる", "なん", "に", "において", "における", "について",
		"にて", "によって", "により", "による", "に対して", "に対する", "に関する", "の", "ので", "のみ", "は", "ば", "へ", "ほか",
		"ほとんど", "ほど", "ます", "また", "または", "まで", "も", "もの", "ものの", "や", "よう", "より", "ら", "られ", "られる", "れ",
		"れる", "を", "ん", "何", "及び", "彼", "彼女", "我々", "特に", "私", "私達", "貴方", "貴方方",
	}
	stopWords := make(map[string]struct{})
	for _, stopWord := range stopWordsSrc {
		stopWords[stopWord] = struct{}{}
	}
	return &japaneseStopwords{
		stopWords: stopWords,
	}
}

func (f japaneseStopwords) Filter(stream token.TokenStream) token.TokenStream {
	result := make(token.TokenStream, 0)
	for _, v := range stream {
		t := v.Term()
		if t.TermType() != term.Text {
			result = append(result, v)
			continue
		}
		if _, ok := f.stopWords[t.Text()]; !ok {
			result = append(result, v)
		}
	}
	return result
}
