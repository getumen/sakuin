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
