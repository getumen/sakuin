package charfilter

import "golang.org/x/text/unicode/norm"

type unicodeNFKCFilter struct {
}

func NewUnicodeNFKCFilter() *unicodeNFKCFilter {
	return &unicodeNFKCFilter{}
}

func (f unicodeNFKCFilter) Filter(s string) string {
	return norm.NFKC.String(s)
}
