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

func (p PostingListCursor) Current() *posting.Posting {
	return p.postingList[p.cursor]
}

func (p *PostingListCursor) Skip(docID int64) {
	p.cursor += posting.ExponentialSearch(p.postingList[p.cursor:], docID)
}

func (p *PostingListCursor) Next() *posting.Posting {
	result := p.postingList[p.cursor]
	p.cursor++
	return result
}

func (p PostingListCursor) Valid() bool {
	return p.cursor < len(p.postingList)
}
