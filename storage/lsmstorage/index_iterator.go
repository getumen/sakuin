package lsmstorage

import (
	"github.com/getumen/sakuin/invertedindex"
)

type indexIterator struct {
	indexes []*invertedindex.InvertedIndex
	cur     int
}

func NewIndexIterator(
	indexes []*invertedindex.InvertedIndex,
) *indexIterator {
	return &indexIterator{
		indexes: indexes,
		cur:     0,
	}
}

func (it indexIterator) HasNext() bool {
	return it.cur < len(it.indexes)
}

func (it *indexIterator) Next() *invertedindex.InvertedIndex {
	result := it.indexes[it.cur]
	it.cur++
	return result
}
