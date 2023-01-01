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
	GetIndex(
		ctx context.Context,
		conds []*termcond.TermCondition,
	) (*invertedindex.InvertedIndex, error)
}
