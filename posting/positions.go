package posting

type positions []uint32

func (p *positions) merge(other positions) {
	newPositions := make([]uint32, 0)
	var left, right int
	for left < len(*p) && right < len(other) {
		if (*p)[left] < other[right] {
			newPositions = append(newPositions, (*p)[left])
			left++
		} else if (*p)[left] == other[right] {
			newPositions = append(newPositions, (*p)[left])
			left++
			right++
		} else {
			newPositions = append(newPositions, other[right])
			right++
		}
	}
	for ; left < len(*p); left++ {
		newPositions = append(newPositions, (*p)[left])
	}
	for ; right < len(other); right++ {
		newPositions = append(newPositions, other[right])
	}
	*p = newPositions
}

func (p positions) Copy() positions {
	v := make([]uint32, len(p))
	copy(v, p)
	return v
}

func (p positions) Cursor() *positionsCursor {
	return &positionsCursor{
		positions: p,
		cur:       0,
	}
}
