package storage

import (
	"context"

	"github.com/getumen/sakuin/invertedindex"
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
	) (IndexIterator, error)
	Close() error
}

type IndexIterator interface {
	HasNext() bool
	Next() (*invertedindex.InvertedIndex, error)
}
