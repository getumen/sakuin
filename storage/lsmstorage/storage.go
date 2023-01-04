package lsmstorage

import (
	"context"
	"fmt"

	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/invertedindex"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type lsmStorage struct {
	index *leveldb.DB
}

func NewLSMStorage(path string) (*lsmStorage, error) {
	db, err := leveldb.OpenFile(
		path,
		&opt.Options{
			Filter: filter.NewBloomFilter(10),
		},
	)
	if err != nil {
		return nil, err
	}
	return &lsmStorage{
		index: db,
	}, nil
}

func (m *lsmStorage) Merge(
	ctx context.Context,
	index *invertedindex.InvertedIndex,
) error {
	tx, err := m.index.OpenTransaction()
	if err != nil {
		return fmt.Errorf("fail to open tx: %w", err)
	}
	it := index.Iterator()

	batch := new(leveldb.Batch)

	for it.Next() {
		key := it.Key().(term.Term)
		value := it.Value().(fieldindex.FieldIndex)
		diskIndexBlob, err := tx.Get(key.Raw(), nil)
		if err == leveldb.ErrNotFound {
			memoryIndexBlob, err := value.Serialize()
			if err != nil {
				return fmt.Errorf("serialize error: %w", err)
			}
			batch.Put(key, memoryIndexBlob)
		} else if err != nil {
			return fmt.Errorf("fail to get index: %w", err)
		} else {
			diskIndex, err := fieldindex.Deserialize(diskIndexBlob)
			if err != nil {
				return fmt.Errorf("deserialize error: %w", err)
			}
			diskIndex.Merge(value)
			diskIndexBlob, err = diskIndex.Serialize()
			if err != nil {
				return fmt.Errorf("serialize error: %w", err)
			}
			batch.Put(key, diskIndexBlob)
		}
	}

	tx.Write(batch, nil)

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	return nil
}

func (m *lsmStorage) GetIndex(
	ctx context.Context,
	conds []*termcond.TermCondition,
) (*invertedindex.InvertedIndex, error) {
	snapshot, err := m.index.GetSnapshot()
	if err != nil {
		return nil, fmt.Errorf("fail to get snapshot: %w", err)
	}

	result := invertedindex.NewInvertedIndex()

	for _, cond := range conds {
		if cond.IsEqual() {
			blob, err := snapshot.Get(cond.Start().Raw(), nil)
			if err == leveldb.ErrNotFound {
				result.Put(cond.Start(), fieldindex.NewFieldIndex())
			} else if err != nil {
				return nil, fmt.Errorf("fail to get index: %w", err)
			} else {
				fieldIndex, err := fieldindex.Deserialize(blob)
				if err != nil {
					return nil, fmt.Errorf("fail to deserialize index: %w", err)
				}
				result.Put(cond.Start(), fieldIndex)
			}
		} else {
			iter := snapshot.NewIterator(nil, nil)
			for ok := iter.Seek(cond.Start().Raw()); ok; ok = iter.Next() {
				key := iter.Key()
				if (!cond.IncludeEnd() && term.Comparator(cond.End(), term.Term(key)) == 0) ||
					term.Comparator(cond.End(), term.Term(key)) < 0 {
					break
				}
				if !cond.IncludeStart() && term.Comparator(cond.Start(), term.Term(key)) == 0 {
					continue
				}
				// cond contains key
				fieldIndex, err := fieldindex.Deserialize(iter.Value())
				if err != nil {
					return nil, fmt.Errorf("fail to deserialize index: %w", err)
				}
				result.Put(term.Term(key), fieldIndex)
			}
			iter.Release()
			if err := iter.Error(); err != nil {
				return nil, fmt.Errorf("iter error: %w", err)
			}
		}
	}

	return result, nil
}

func (m *lsmStorage) Close() error {
	return m.index.Close()
}
