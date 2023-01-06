package posting

import (
	"encoding/binary"
	"fmt"

	"github.com/getumen/sakuin/encoding/int32enc"
	"github.com/getumen/sakuin/encoding/int64enc"
)

type Posting struct {
	docID     uint64
	positions positions
}

func NewPosting(
	docID uint64,
	positions []uint32,
) *Posting {
	return &Posting{
		docID:     docID,
		positions: positions,
	}
}

func (p Posting) DocIDForTest() uint64 {
	return p.docID
}

func (p Posting) Compare(other *Posting) int {
	return int(p.docID) - int(other.docID)
}

func (p *Posting) Max(other *Posting) *Posting {
	if p.docID < other.docID {
		return other
	}
	return p
}

func (p Posting) Copy() *Posting {
	return &Posting{
		docID:     p.docID,
		positions: p.positions.Copy(),
	}
}

func (p *Posting) Merge(other *Posting) {
	p.positions.merge(other.positions)
}

func (p Posting) EstimateSize() int {
	return len(p.positions) * 4
}

func Serialize(postings []*Posting) ([]byte, error) {
	result := make([]byte, 0)

	docIDs := make([]uint64, len(postings))
	for i := range postings {
		docIDs[i] = postings[i].docID
	}
	int64enc.EncodeDeltaInplace(docIDs)

	docIDBlob, err := int64enc.Encode(docIDs)
	if err != nil {
		return nil, fmt.Errorf("encode doc id error: %w", err)
	}
	result = binary.AppendUvarint(result, uint64(len(docIDBlob)))
	result = append(result, docIDBlob...)

	positionList := make([]uint32, 0)
	for i := range postings {
		ps := make([]uint32, len(postings[i].positions))
		copy(ps, postings[i].positions)
		int32enc.EncodeDeltaInplace(ps)
		positionList = append(positionList, uint32(len(postings[i].positions)))
		positionList = append(positionList, ps...)
	}
	positionBlob, err := int32enc.Encode(positionList)
	if err != nil {
		return nil, fmt.Errorf("encode position error: %w", err)
	}
	result = binary.AppendUvarint(result, uint64(len(positionBlob)))
	result = append(result, positionBlob...)

	return result, nil
}

func Deserialize(blob []byte) ([]*Posting, error) {
	var cur int
	docIDBlobSize, diff := binary.Uvarint(blob[cur:])
	cur += diff
	docIDs, err := int64enc.Decode(blob[cur : cur+int(docIDBlobSize)])
	if err != nil {
		return nil, fmt.Errorf("decode doc id failed: %w", err)
	}
	cur += int(docIDBlobSize)
	int64enc.DecodeDeltaInplace(docIDs)

	positionBlobSize, diff := binary.Uvarint(blob[cur:])
	cur += diff
	positionList, err := int32enc.Decode(blob[cur : cur+int(positionBlobSize)])
	if err != nil {
		return nil, fmt.Errorf("decode positions failed: %w", err)
	}

	var posIndex int

	result := make([]*Posting, len(docIDs))
	for i, docID := range docIDs {
		n := positionList[posIndex]
		posIndex++
		ps := make([]uint32, n)
		for p := 0; p < int(n); p++ {
			ps[p] = positionList[posIndex]
			posIndex++
		}
		int32enc.DecodeDeltaInplace(ps)

		result[i] = NewPosting(docID, ps)
	}

	return result, nil
}

func PhraseMatch(postingLists []*Posting, relativePosition []uint32) *Posting {
	// positionがマッチするかを探索
	positionCursors := make([]*positionsCursor, len(postingLists))
	for i, v := range postingLists {
		if len(v.positions) == 0 {
			return nil
		}
		positionCursors[i] = v.positions.Cursor()
	}

	matchPositions := make([]uint32, 0)

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
		return NewPosting(postingLists[0].docID, matchPositions)
	}
	return nil
}
