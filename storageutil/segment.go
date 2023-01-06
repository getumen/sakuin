package storageutil

import (
	"encoding/binary"
	"fmt"

	"github.com/getumen/sakuin/fieldindex"
)

type Segment struct {
	blob []byte
}

func NewSegment(
	blob []byte,
) *Segment {
	return &Segment{
		blob: blob,
	}
}

func (s *Segment) Iterator() *SegmentIterator {
	return &SegmentIterator{
		seg: s,
		cur: 0,
	}
}

func (s *Segment) FindAvailableSegment(newIndexSize, maxSegmentSize int) int {
	var cur int
	segmentID := 1
	for cur < len(s.blob) {
		size, n := binary.Uvarint(s.blob[cur:])
		cur += n + int(size)
		if int(size)+newIndexSize < maxSegmentSize {
			return segmentID
		}
		segmentID++
	}
	return segmentID
}

func (s *Segment) Save(segmentID int, fieldIndex fieldindex.FieldIndex) error {
	newBlob := make([]byte, 0)
	var cur int

	blob, err := fieldIndex.Serialize()
	if err != nil {
		return fmt.Errorf("failt to save segment: %w", err)
	}

	for i := 1; i < segmentID; i++ {
		if cur < len(s.blob) {
			size, n := binary.Uvarint(s.blob[cur:])
			cur += n
			newBlob = binary.AppendUvarint(newBlob, size)
			newBlob = append(newBlob, s.blob[cur:cur+int(size)]...)
			cur += int(size)
		} else {
			newBlob = binary.AppendUvarint(newBlob, 0)
		}
	}

	// currentSegmentID == segmentID
	newBlob = binary.AppendUvarint(newBlob, uint64(len(blob)))
	newBlob = append(newBlob, blob...)
	size, n := binary.Uvarint(s.blob[cur:])
	cur += n + int(size)

	// currentSegmentID > segmentID
	for cur < len(s.blob) {
		size, n := binary.Uvarint(s.blob[cur:])
		cur += n
		newBlob = binary.AppendUvarint(newBlob, size)
		newBlob = append(newBlob, s.blob[cur:cur+int(size)]...)
		cur += int(size)
	}
	s.blob = newBlob
	return nil
}

func (s Segment) Get(segmentID int) (index fieldindex.FieldIndex, err error) {
	var cur int
	currentSegmentID := 1
	for cur < len(s.blob) {
		if currentSegmentID < segmentID {
			size, n := binary.Uvarint(s.blob[cur:])
			cur += n + int(size)
			currentSegmentID++
		} else if currentSegmentID == segmentID {
			size, n := binary.Uvarint(s.blob[cur:])
			cur += n
			return fieldindex.Deserialize(s.blob[cur : cur+int(size)])
		} else {
			return fieldindex.NewFieldIndex(), nil
		}
	}
	return fieldindex.NewFieldIndex(), nil
}

func (s Segment) Bytes() []byte {
	return s.blob
}
