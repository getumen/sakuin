package postinglist

import "github.com/getumen/sakuin/posting"

type PostingListCursor struct {
	PostingList
	cursor int
}

func NewPostingListCursor(p PostingList) *PostingListCursor {
	return &PostingListCursor{
		PostingList: p,
		cursor:      0,
	}
}

func (p PostingListCursor) Value() *posting.Posting {
	return p.postingList[p.cursor]
}

func (p *PostingListCursor) Skip(otherPosting *posting.Posting) bool {
	p.cursor += posting.ExponentialSearch(
		p.postingList[p.cursor:],
		otherPosting,
	)
	return p.cursor < p.Len()
}

func (p *PostingListCursor) Next() bool {
	p.cursor++
	return p.cursor < p.Len()
}

func (p *PostingListCursor) Compare(other *PostingListCursor) int {
	return p.Value().Compare(other.Value())
}
