package fieldindex

import (
	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/postinglist"
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
