package posting

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func ExponentialSearch(arr []*Posting, posting *Posting) int {
	if arr[0].docID == posting.docID {
		return 0
	}
	cur := 1
	for cur < len(arr) && arr[cur].docID <= posting.docID {
		cur = cur * 2
	}
	return cur/2 + BinarySearch(
		arr[cur/2:min(len(arr), cur)],
		posting,
	)
}

func BinarySearch(arr []*Posting, posting *Posting) int {
	left, right := 0, len(arr)
	for left < right {
		h := int(uint(left+right) >> 1)
		if arr[h].docID < posting.docID {
			left = h + 1
		} else {
			right = h
		}
	}
	return left
}
