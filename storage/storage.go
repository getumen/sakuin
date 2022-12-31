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
	) (err error)
	GetIndex(
		ctx context.Context,
		terms []*termcond.TermCondition,
	) (index *invertedindex.InvertedIndex, err error)
}
