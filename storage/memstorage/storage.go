package memstorage

import (
	"context"
	"sync"

	"github.com/getumen/sakuin/invertedindex"
	"github.com/getumen/sakuin/termcond"
)

type memStorage struct {
	mu    sync.RWMutex
	index *invertedindex.InvertedIndex
}

func NewMemStorage() *memStorage {
	return &memStorage{
		mu:    sync.RWMutex{},
		index: invertedindex.NewInvertedIndex(),
	}
}

func (m *memStorage) Merge(
	ctx context.Context,
	index *invertedindex.InvertedIndex,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.index.Merge(index)
	return nil
}

func (m *memStorage) GetIndex(
	ctx context.Context,
	conds []*termcond.TermCondition,
) (index *invertedindex.InvertedIndex, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.index.GetPartialIndex(conds), nil
}

func (m *memStorage) Close() error {
	m.index = nil
	return nil
}
