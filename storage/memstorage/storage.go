package memstorage

import (
	"context"

	"github.com/getumen/sakuin/invertedindex"
	"github.com/getumen/sakuin/termcond"
)

type memStorage struct {
	index *invertedindex.InvertedIndex
}

func NewMemStorage() *memStorage {
	return &memStorage{
		index: invertedindex.NewInvertedIndex(),
	}
}

func (m *memStorage) Merge(
	ctx context.Context,
	index *invertedindex.InvertedIndex,
) error {
	m.index.Merge(index)
	return nil
}

func (m *memStorage) GetIndex(
	ctx context.Context,
	conds []*termcond.TermCondition,
) (index *invertedindex.InvertedIndex, err error) {
	return m.index.GetPartialIndex(conds), nil
}
