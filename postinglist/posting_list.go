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

func Intersection(postingLists []*PostingList, relativePosition []int64) *PostingList {

	if len(postingLists) == 0 {
		return NewPostingList(make([]*posting.Posting, 0))
	}

	postingListCursors := make([]*PostingListCursor, len(postingLists))
	var maxDocID int64 = 0
	for i := range postingLists {
		if postingLists[i].Len() == 0 {
			return NewPostingList(make([]*posting.Posting, 0))
		}
		postingListCursors[i] = postingLists[i].Cursor()
		maxDocID = max(maxDocID, postingListCursors[i].Current().GetDocID())
	}

	result := make([]*posting.Posting, 0)

POSTING_LOOP:
	for {
		matchCount := 0
		for postingListCursorIndex := range postingListCursors {

			if postingListCursors[postingListCursorIndex].Current().GetDocID() < maxDocID {
				postingListCursors[postingListCursorIndex].Skip(maxDocID)
				// カーソルを読み終わったら終了
				if !postingListCursors[postingListCursorIndex].Valid() {
					break POSTING_LOOP
				}
				// カーソルが読み終わっていなかったらmaxDocIDを更新して次のmatchCountのカウントへ
				maxDocID = max(maxDocID, postingListCursors[postingListCursorIndex].Current().GetDocID())
				break
			}
			matchCount++
		}
		if matchCount < len(postingListCursors) {
			continue
		}
		// 空の場合はフレーズマッチではないので、DocIDの積を返す
		if relativePosition == nil {
			result = append(result, posting.NewPosting(maxDocID, nil))
			// 次のドキュメントを探す
			maxDocID++
			continue
		}
		positionCursors := make([]*position.PositionsCursor, len(postingListCursors))
		for i, v := range postingListCursors {
			positionCursors[i] = v.Current().GetPositions().Cursor()
		}

	POSITION_LOOP:
		for {
			positionMatchCount := 0
			var currentOffset int64
			for positionCursorIndex := range positionCursors {
				if positionCursorIndex == 0 {
					positionMatchCount++
					currentOffset = positionCursors[positionCursorIndex].Current()
					continue
				}
				for positionCursors[positionCursorIndex].Current() < currentOffset+relativePosition[positionCursorIndex] {
					positionCursors[positionCursorIndex].Skip(currentOffset + relativePosition[positionCursorIndex])
					if !positionCursors[positionCursorIndex].Valid() {
						break POSITION_LOOP
					}
				}
				if positionCursors[positionCursorIndex].Current() == currentOffset+relativePosition[positionCursorIndex] {
					positionMatchCount++
				}
			}
			if positionMatchCount == len(positionCursors) {
				result = append(result, posting.NewPosting(maxDocID, nil))
				// 次のドキュメントを探す
				maxDocID++
				break POSITION_LOOP
			}
		}
	}

	return NewPostingList(result)
}

type unionHeap []*PostingListCursor

func (h unionHeap) Len() int { return len(h) }

func (h unionHeap) Less(i, j int) bool {
	return h[i].Current().GetDocID() < h[j].Current().GetDocID()
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
		cursor := postingLists[i].Cursor()
		if cursor.Valid() {
			postingListCursors = append(postingListCursors, cursor)
		}
	}
	heap.Init(&postingListCursors)

	result := make([]*posting.Posting, 0)

	var minDocID int64

	for postingListCursors.Len() > 0 {
		cursor := heap.Pop(&postingListCursors).(*PostingListCursor)
		if cursor.Current().GetDocID() != minDocID {
			minDocID = cursor.Current().GetDocID()
			result = append(result, posting.NewPosting(minDocID, nil))
		}
		cursor.Next()
		if cursor.Valid() {
			heap.Push(&postingListCursors, cursor)
		}
	}
	return NewPostingList(result)
}

func Difference(x, y *PostingList) *PostingList {
	xCursor := x.Cursor()
	yCursor := y.Cursor()

	result := make([]*posting.Posting, 0)
	for xCursor.Valid() {
		if !yCursor.Valid() {
			for xCursor.Valid() {
				result = append(result, xCursor.Next())
			}
			break
		}
		xPosting := xCursor.Next()
		yCursor.Skip(xPosting.GetDocID())

		if yCursor.Valid() && yCursor.Current().GetDocID() == xPosting.GetDocID() {
			continue
		}
		result = append(result, xPosting)
	}

	return NewPostingList(result)
}
