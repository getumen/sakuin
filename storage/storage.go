package storage

import (
	"context"
	"fmt"

	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/invertedindex"
	"github.com/getumen/sakuin/storageutil"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
)

type IndexStorage interface {
	Merge(
		ctx context.Context,
		index *invertedindex.InvertedIndex,
	) error
	GetIndexIterator(
		ctx context.Context,
		conds []*termcond.TermCondition,
	) (*IndexIterator, error)
	Close() error
}

type IndexIterator struct {
	terms            []term.Term
	segmentIterators []*storageutil.SegmentIterator
	segmentID        int
	finishedSegment  []bool
}

func NewIndexIterator(
	terms []term.Term,
	segmentIterators []*storageutil.SegmentIterator,
) *IndexIterator {
	return &IndexIterator{
		terms:            terms,
		segmentIterators: segmentIterators,
		segmentID:        1,
		finishedSegment:  make([]bool, len(terms)),
	}
}

func (it IndexIterator) HasNext() bool {
	for i := range it.terms {
		if !it.finishedSegment[i] {
			return true
		}
	}
	return false
}

func (it *IndexIterator) Next() (*invertedindex.InvertedIndex, error) {
	index := invertedindex.NewInvertedIndex(it.segmentID)
	it.segmentID++
	for i := range it.terms {
		if it.segmentIterators[i].HasNext() {
			fieldIndex, err := it.segmentIterators[i].Next()
			if err != nil {
				return nil, fmt.Errorf("fail to deserialize: %w", err)
			}
			index.Put(it.terms[i], fieldIndex)
			it.finishedSegment[i] = !it.segmentIterators[i].HasNext()
		} else {
			index.Put(it.terms[i], fieldindex.NewFieldIndex())
		}
	}
	return index, nil
}
