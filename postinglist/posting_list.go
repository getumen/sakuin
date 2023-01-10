package postinglist

import (
	"container/heap"
	"fmt"

	"github.com/getumen/sakuin/posting"
)

type PostingList struct {
	postingList []*posting.Posting
}

func NewPostingList(
	pl []*posting.Posting,
) *PostingList {
	return &PostingList{
		postingList: pl,
	}
}

func (p *PostingList) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p PostingList) Serialize() ([]byte, error) {
	return posting.Serialize(p.postingList)
}

func Deserialize(blob []byte) (*PostingList, error) {
	pl, err := posting.Deserialize(blob)
	if err != nil {
		return nil, fmt.Errorf("deserialize error: %w", err)
	}
	return &PostingList{
		postingList: pl,
	}, nil
}

func (p PostingList) EstimateSize() int {
	var size int
	for _, posting := range p.postingList {
		size += posting.EstimateSize()
	}
	return size
}

func (p *PostingList) Merge(other *PostingList) {
	newPostings := make([]*posting.Posting, 0)
	var left, right int
	for left < len(p.postingList) && right < len(other.postingList) {
		if p.postingList[left].Compare(other.postingList[right]) < 0 {
			newPostings = append(newPostings, p.postingList[left])
			left++
		} else if p.postingList[left].Compare(other.postingList[right]) == 0 {
			p.postingList[left].Merge(other.postingList[right])
			newPostings = append(newPostings, p.postingList[left])
			left++
			right++
		} else {
			newPostings = append(newPostings, other.postingList[right])
			right++
		}
	}
	for ; left < len(p.postingList); left++ {
		newPostings = append(newPostings, p.postingList[left])
	}
	for ; right < len(other.postingList); right++ {
		newPostings = append(newPostings, other.postingList[right])
	}
	p.postingList = newPostings
}

func (p PostingList) Cursor() *PostingListCursor {
	return NewPostingListCursor(p)
}

func (p PostingList) Len() int {
	return len(p.postingList)
}

func PhraseMatch(postingLists []*PostingList, relativePosition []uint32) *PostingList {
	if len(postingLists) == 0 {
		return NewPostingList(make([]*posting.Posting, 0))
	}

	postingListCursors := make([]*PostingListCursor, len(postingLists))

	for i := range postingLists {
		if postingLists[i].Len() == 0 {
			return NewPostingList(make([]*posting.Posting, 0))
		}
		postingListCursors[i] = postingLists[i].Cursor()
	}

	result := make([]*posting.Posting, 0)

POSTING_LOOP:
	for searchNextMatchDocID(postingListCursors) {

		postingList := make([]*posting.Posting, len(postingListCursors))
		for i := range postingListCursors {
			postingList[i] = postingListCursors[i].Value()
		}

		matchPosting := posting.PhraseMatch(postingList, relativePosition)
		if matchPosting != nil {
			result = append(result, matchPosting)
		}

		for index := range postingListCursors {
			cursor := postingListCursors[index]
			if !cursor.Next() {
				break POSTING_LOOP
			}
		}
	}

	return NewPostingList(result)
}

func searchNextMatchDocID(postingListCursors []*PostingListCursor) bool {
	maxPosting := posting.NewPosting(0, nil)
	for {
		matchCount := 0
		// すべてのカーソルに存在するdocIDを探す
		for index := range postingListCursors {
			cursor := postingListCursors[index]
			if cursor.Value().Compare(maxPosting) < 0 {
				// カーソルを読み終わったら終了
				if !cursor.Skip(maxPosting) {
					return false
				}
				maxPosting = maxPosting.Max(cursor.Value())
				break
			} else if cursor.Value().Compare(maxPosting) > 0 {
				maxPosting = cursor.Value()
			} else {
				matchCount++
			}
		}
		if matchCount < len(postingListCursors) {
			return true
		}
	}
}

func Intersection(postingLists []*PostingList) *PostingList {
	if len(postingLists) == 0 {
		return NewPostingList(make([]*posting.Posting, 0))
	}

	postingListCursors := make([]*PostingListCursor, len(postingLists))

	for i := range postingLists {
		if postingLists[i].Len() == 0 {
			return NewPostingList(make([]*posting.Posting, 0))
		}
		postingListCursors[i] = postingLists[i].Cursor()
	}

	result := make([]*posting.Posting, 0)

POSTING_LOOP:
	for searchNextMatchDocID(postingListCursors) {
		// positionをマージ
		var posting *posting.Posting
		for _, v := range postingListCursors {
			if posting == nil {
				posting = v.Value().Copy()
				continue
			}
			posting.Merge(v.Value())
		}
		result = append(result, posting)
		for _, v := range postingListCursors {
			if !v.Next() {
				break POSTING_LOOP
			}
		}
	}

	return NewPostingList(result)
}

type unionHeap []*PostingListCursor

func (h unionHeap) Len() int { return len(h) }

func (h unionHeap) Less(i, j int) bool {
	return h[i].Compare(h[j]) < 0
}

func (h unionHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *unionHeap) Push(x any) {
	item := x.(*PostingListCursor)
	*h = append(*h, item)
}

func (h *unionHeap) Pop() any {
	old := *h
	n := old.Len()
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*h = old[0 : n-1]
	return item
}

func Union(postingLists []*PostingList) *PostingList {
	if len(postingLists) == 0 {
		return NewPostingList(make([]*posting.Posting, 0))
	}
	if len(postingLists) == 1 {
		return postingLists[0]
	}

	postingListCursors := make(unionHeap, 0)
	for i := range postingLists {
		if postingLists[i].Len() > 0 {
			postingListCursors = append(postingListCursors, postingLists[i].Cursor())
		}
	}
	heap.Init(&postingListCursors)

	result := make([]*posting.Posting, 0)

	minPosting := posting.NewPosting(0, nil)

	for postingListCursors.Len() > 0 {
		cursor := heap.Pop(&postingListCursors).(*PostingListCursor)
		if cursor.Value().Compare(minPosting) != 0 {
			minPosting = cursor.Value().Copy()
			result = append(result, minPosting)
		} else {
			// merge positions
			result[len(result)-1].Merge(cursor.Value())
		}
		if cursor.Next() {
			heap.Push(&postingListCursors, cursor)
		}
	}
	return NewPostingList(result)
}

func Difference(x, y *PostingList) *PostingList {
	xCursor := x.Cursor()
	yCursor := y.Cursor()

	result := make([]*posting.Posting, 0)

	if xCursor.Len() == 0 {
		return NewPostingList(result)
	}
	if yCursor.Len() == 0 {
		for {
			result = append(result, xCursor.Value())
			if !xCursor.Next() {
				break
			}
		}
		return NewPostingList(result)
	}

LOOP:
	for {
		if xCursor.Compare(yCursor) < 0 {
			result = append(result, x.Cursor().Value())
			if !xCursor.Next() {
				break
			}
			continue
		}
		if xCursor.Compare(yCursor) == 0 {
			if !yCursor.Next() {
				for xCursor.Next() {
					result = append(result, xCursor.Value())
					break LOOP
				}
			}
			if !xCursor.Next() {
				break
			}
		}
		if !yCursor.Skip(xCursor.Value()) {
			for !xCursor.Next() {
				result = append(result, xCursor.Value())
				break LOOP
			}
		}
	}

	return NewPostingList(result)
}
