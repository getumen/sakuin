package position

type PositionsCursor struct {
	Positions
	cur int
}

func (p *PositionsCursor) Skip(value int64) bool {
	p.cur += exponentialSearch(p.value[p.cur:], value)
	return p.cur < p.Len()
}

func (p *PositionsCursor) Next() bool {
	p.cur++
	return p.cur < p.Len()
}

func (p PositionsCursor) Value() int64 {
	return p.value[p.cur]
}

func exponentialSearch(arr []int64, docID int64) int {
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

func binarySearch(arr []int64, docID int64) int {
	i, j := 0, len(arr)
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if arr[h] < docID {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
