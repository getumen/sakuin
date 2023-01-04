package posting

import "sort"

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

var BinarySearch = standardBinarySearchForPosting

func sortBinarySearchForPosting(arr []*Posting, posting *Posting) int {
	return sort.Search(len(arr), func(i int) bool { return arr[i].docID >= posting.docID })
}

func standardBinarySearchForPosting(arr []*Posting, posting *Posting) int {
	left, right := 0, len(arr)
	for left < right {
		h := int(uint(left+right) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if arr[h].docID < posting.docID {
			left = h + 1 // preserves f(i-1) == false
		} else {
			right = h // preserves f(j) == true
		}
	}
	return left
}

func branchPredictionBinarySearchForPosting(arr []*Posting, posting *Posting) int {
	left, right := 0, len(arr)
	for right > 1 {
		half := right / 2
		if arr[left+half-1].docID < posting.docID {
			left += half
			right -= half
		} else {
			right = half
		}
	}
	return left
}

func branchlessBinarySearchForPosting(arr []*Posting, posting *Posting) int {
	left, right := 0, len(arr)
	for right > 1 {
		half := right / 2
		left += boolToInt(arr[left+half-1].docID < posting.docID) * half
		right -= half
	}
	return left
}

func exponentialSearch(arr []uint32, docID uint32) int {
	if arr[0] == docID {
		return 0
	}
	cur := 1
	for cur < len(arr) && arr[cur] <= docID {
		cur = cur * 2
	}
	return cur/2 + binarySearch(
		arr[cur/2:min(len(arr), cur)],
		docID,
	)
}

var binarySearch = standardBinarySearchForPosition

func sortBinarySearchForPosition(arr []uint32, pos uint32) int {
	return sort.Search(len(arr), func(i int) bool { return arr[i] >= pos })
}

func standardBinarySearchForPosition(arr []uint32, pos uint32) int {
	left, right := 0, len(arr)
	for left < right {
		h := int(uint(left+right) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if arr[h] < pos {
			left = h + 1 // preserves f(i-1) == false
		} else {
			right = h // preserves f(j) == true
		}
	}
	return left
}

func branchPredictionBinarySearchForPosition(arr []uint32, pos uint32) int {
	left, right := 0, len(arr)
	for right > 1 {
		half := right / 2
		if arr[left+half-1] < pos {
			left += half
			right -= half
		} else {
			right = half
		}
	}
	return left
}

func branchlessBinarySearchForPosition(arr []uint32, pos uint32) int {
	left, right := 0, len(arr)
	for right > 1 {
		half := right / 2
		left += boolToInt(arr[left+half-1] < pos) * half
		right -= half
	}
	return left
}
func boolToInt(b bool) int {
	// The compiler currently only optimizes this form.
	// See issue 6011.
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}
