package posting

import "github.com/getumen/sakuin/position"

type Posting struct {
	docID     int64
	positions *position.Positions
}

func NewPosting(
	docID int64,
	positions *position.Positions,
) *Posting {
	return &Posting{
		docID:     docID,
		positions: positions,
	}
}

func (p Posting) Compare(other *Posting) int {
	return int(p.docID) - int(other.docID)
}

func (p Posting) GetDocID() int64 {
	return p.docID
}

func (p Posting) GetPositions() *position.Positions {
	return p.positions
}

func (p *Posting) Merge(other *Posting) {
	newPositions := make([]int64, 0)
	var left, right int
	for left < p.positions.Len() && right < other.positions.Len() {
		if p.positions.At(left) < other.positions.At(right) {
			newPositions = append(newPositions, p.positions.At(left))
			left++
		} else if p.positions.At(left) == other.positions.At(right) {
			newPositions = append(newPositions, p.positions.At(left))
			left++
			right++
		} else {
			newPositions = append(newPositions, other.positions.At(right))
			right++
		}
	}
	for ; left < p.positions.Len(); left++ {
		newPositions = append(newPositions, p.positions.At(left))
	}
	for ; right < other.positions.Len(); right++ {
		newPositions = append(newPositions, other.positions.At(right))
	}
	p.positions = position.NewPositions(newPositions)
}
