package postinglist

import (
	"container/heap"

	"github.com/getumen/sakuin/position"
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

func (p PostingList) GetPostingList() []*posting.Posting {
	return p.postingList
}

func (p PostingList) Cursor() *PostingListCursor {
	return NewPostingListCursor(p)
}

func (p PostingList) Len() int {
	return len(p.postingList)
}

func PhraseMatch(postingLists []*PostingList, relativePosition []int64) *PostingList {
	if len(postingLists) == 0 {
		return NewPostingList(make([]*posting.Posting, 0))
	}

	postingListCursors := make([]*PostingListCursor, len(postingLists))
	maxPosting := posting.NewPosting(0, nil)

	for i := range postingLists {
		if postingLists[i].Len() == 0 {
			return NewPostingList(make([]*posting.Posting, 0))
		}
		postingListCursors[i] = postingLists[i].Cursor()
		maxPosting = maxPosting.Max(postingListCursors[i].Value())
	}

	result := make([]*posting.Posting, 0)

POSTING_LOOP:
	for {
		matchCount := 0
		// すべてのカーソルに存在するdocIDを探す
		for index := range postingListCursors {
			cursor := postingListCursors[index]

			if cursor.Value().GetDocID() < maxPosting.GetDocID() {
				// カーソルを読み終わったら終了
				if !cursor.Skip(maxPosting) {
					break POSTING_LOOP
				}
				maxPosting = maxPosting.Max(cursor.Value())
				break
			}
			matchCount++
		}
		if matchCount < len(postingListCursors) {
			continue
		}

		// positionがマッチするかを探索
		positionCursors := make([]*position.PositionsCursor, len(postingListCursors))
		for i, v := range postingListCursors {
			positionCursors[i] = v.Value().GetPositions().Cursor()
		}

		if positionCursors[0].Len() == 0 {
			for index := range postingListCursors {
				if !postingListCursors[index].Next() {
					break POSTING_LOOP
				}
			}
		}

		matchPositions := make([]int64, 0)

	POSITION_LOOP:
		for {
			positionMatchCount := 1
			currentOffset := positionCursors[0].Value()
			for index := 1; index < len(positionCursors); index++ {
				absolutePosition := currentOffset + relativePosition[index]
				for positionCursors[index].Value() < absolutePosition {
					if !positionCursors[index].Skip(absolutePosition) {
						break POSITION_LOOP
					}
				}
				if positionCursors[index].Value() == absolutePosition {
					positionMatchCount++
				}
			}

			if positionMatchCount == len(positionCursors) {
				matchPositions = append(matchPositions, positionCursors[0].Value())
			}
			if !positionCursors[0].Next() {
				break POSITION_LOOP
			}
		}
		if len(matchPositions) > 0 {
			result = append(
				result,
				posting.NewPosting(
					postingListCursors[0].Value().GetDocID(),
					position.NewPositions(matchPositions),
				),
			)
		}

		for index := range postingListCursors {
			cursor := postingListCursors[index]
			if !cursor.Next() {
				break POSTING_LOOP
			}
			maxPosting = maxPosting.Max(cursor.Value())
		}
	}

	return NewPostingList(result)
}

func Intersection(postingLists []*PostingList) *PostingList {
	if len(postingLists) == 0 {
		return NewPostingList(make([]*posting.Posting, 0))
	}

	postingListCursors := make([]*PostingListCursor, len(postingLists))
	maxPosting := posting.NewPosting(0, nil)
	for i := range postingLists {
		if postingLists[i].Len() == 0 {
			return NewPostingList(make([]*posting.Posting, 0))
		}
		postingListCursors[i] = postingLists[i].Cursor()
		maxPosting = maxPosting.Max(postingListCursors[i].Value())
	}

	result := make([]*posting.Posting, 0)

POSTING_LOOP:
	for {
		matchCount := 0
		// すべてのカーソルに存在するdocIDを探す
		for index := range postingListCursors {
			if postingListCursors[index].Value().Compare(maxPosting) < 0 {
				// カーソルを読み終わったら終了
				if !postingListCursors[index].Skip(maxPosting) {
					break POSTING_LOOP
				}
				maxPosting = maxPosting.Max(postingListCursors[index].Value())
				break
			}
			matchCount++
		}
		if matchCount < len(postingListCursors) {
			continue
		}

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
