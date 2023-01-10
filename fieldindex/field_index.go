package fieldindex

import (
	"fmt"

	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/postinglist"
	"google.golang.org/protobuf/proto"
)

type FieldIndex map[fieldname.FieldName]*postinglist.PostingList

func NewFieldIndex() FieldIndex {
	return make(FieldIndex)
}

func NewFieldIndexFromMap(
	value map[fieldname.FieldName]*postinglist.PostingList,
) FieldIndex {
	return value
}

func (f FieldIndex) Merge(other FieldIndex) {
	for key, value := range other {
		if _, ok := f[key]; ok {
			f[key].Merge(value)
		} else {
			f[key] = value
		}
	}
}

func (f FieldIndex) EstimateSize() int {
	var size int
	for k, v := range f {
		size += len(k)
		size += v.EstimateSize()
	}
	return size
}

func Deserialize(blob []byte) (FieldIndex, error) {
	p := &Record{}

	if err := proto.Unmarshal(blob, p); err != nil {
		return nil, fmt.Errorf("unmarshal field index error: %w", err)
	}
	result := make(FieldIndex)
	for key, value := range p.GetFieldIndex() {
		blob, err := postinglist.Deserialize(value)
		if err != nil {
			return nil, fmt.Errorf("deserialize index error: %w", err)
		}
		result[fieldname.FieldName(key)] = blob
	}
	return result, nil
}

func (f FieldIndex) Serialize() ([]byte, error) {
	index := make(map[string][]byte)
	for fieldName, postingList := range f {
		pl, err := postingList.Serialize()
		if err != nil {
			return nil, fmt.Errorf("serialize index error: %w", err)
		}
		index[string(fieldName)] = pl
	}
	p := Record{FieldIndex: index}
	if blob, err := proto.Marshal(&p); err != nil {
		return nil, fmt.Errorf("marshal field index error: %w", err)
	} else {
		return blob, nil
	}
}
