package posting

type positionsCursor struct {
	positions
	cur int
}

func (p *positionsCursor) Skip(value uint32) bool {
	p.cur += exponentialSearch(p.positions[p.cur:], value)
	return p.cur < len(p.positions)
}

func (p *positionsCursor) Next() bool {
	p.cur++
	return p.cur < len(p.positions)
}

func (p positionsCursor) Value() uint32 {
	return p.positions[p.cur]
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

func binarySearch(arr []uint32, docID uint32) int {
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
