package storageutil

import (
	"encoding/binary"
	"fmt"

	"github.com/getumen/sakuin/fieldindex"
)

type SegmentIterator struct {
	seg *Segment
	cur int
}

func (s *SegmentIterator) HasNext() bool {
	return s.cur < len(s.seg.blob)
}

func (s *SegmentIterator) Next() (fieldindex.FieldIndex, error) {
	size, n := binary.Uvarint(s.seg.blob[s.cur:])
	s.cur += n
	if size == 0 {
		return fieldindex.NewFieldIndex(), nil
	} else {
		result, err := fieldindex.Deserialize(s.seg.blob[s.cur : s.cur+int(size)])
		if err != nil {
			return nil, fmt.Errorf("fail to deserialize field index: %w", err)
		}
		s.cur += int(size)
		return result, nil
	}
}
