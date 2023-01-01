package position

type Positions struct {
	value []int64
}

func NewPositions(value []int64) *Positions {
	return &Positions{
		value: value,
	}
}

func (p Positions) Len() int {
	return len(p.value)
}

func (p Positions) Copy() *Positions {
	v := make([]int64, len(p.value))
	copy(v, p.value)
	return &Positions{
		value: v,
	}
}

func (p Positions) At(i int) int64 {
	return p.value[i]
}

func (p Positions) Cursor() *PositionsCursor {
	return &PositionsCursor{
		Positions: p,
		cur:       0,
	}
}
