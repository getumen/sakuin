package lsmstorage

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/invertedindex"
	"github.com/getumen/sakuin/storage"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type lsmStorage struct {
	index *leveldb.DB
}

const maxSegmentSize = 1024 * 1024

func NewStorage(path string) (*lsmStorage, error) {
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

func getTerm(key []byte) term.Term {
	result := make([]byte, len(key)-4)
	copy(result, key[:len(key)-4])
	return result
}

func getSegmentID(key []byte) uint32 {
	return binary.LittleEndian.Uint32(key[len(key)-4:])
}

func getDBKey(term term.Term, segmentID uint32) []byte {
	key := make([]byte, 0)
	key = append(key, term...)
	key = binary.LittleEndian.AppendUint32(key, segmentID)
	return key
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

	availableSegment := make([]bool, 0)

	for it.Next() {
		key := it.Key().(term.Term)
		value := it.Value().(fieldindex.FieldIndex)

		dbIter := tx.NewIterator(util.BytesPrefix(key.Raw()), nil)

		for dbIter.Next() {
			segmentID := getSegmentID(dbIter.Key())

			for i := len(availableSegment); i <= int(segmentID); i++ {
				availableSegment = append(availableSegment, true)
			}

			availableSegment[segmentID] = availableSegment[segmentID] && len(dbIter.Value())+value.EstimateSize() < maxSegmentSize
		}
	}

	it = index.Iterator()
	batch := new(leveldb.Batch)

	availableSegmentID := len(availableSegment)
	for i := 0; i < len(availableSegment); i++ {
		if availableSegment[i] {
			availableSegmentID = i
			break
		}
	}

	for it.Next() {
		dbKey := getDBKey(it.Key().(term.Term), uint32(availableSegmentID))
		memoryIndex := it.Value().(fieldindex.FieldIndex)
		value, err := tx.Get(dbKey, nil)
		if errors.Is(err, leveldb.ErrNotFound) {
			blob, err := memoryIndex.Serialize()
			if err != nil {
				return fmt.Errorf("fail to merge: %w", err)
			}
			batch.Put(dbKey, blob)
			continue
		}
		if err != nil {
			return fmt.Errorf("fail to merge: %w", err)
		}
		dbIndex, err := fieldindex.Deserialize(value)
		if err != nil {
			return fmt.Errorf("fail to merge: %w", err)
		}
		dbIndex.Merge(memoryIndex)
		blob, err := dbIndex.Serialize()
		if err != nil {
			return fmt.Errorf("fail to merge: %w", err)
		}
		batch.Put(dbKey, blob)
	}

	err = tx.Write(batch, nil)
	if err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	return nil
}

func (m *lsmStorage) GetIndexIterator(
	ctx context.Context,
	conds []*termcond.TermCondition,
) (storage.IndexIterator, error) {
	snapshot, err := m.index.GetSnapshot()
	if err != nil {
		return nil, fmt.Errorf("fail to get snapshot: %w", err)
	}

	indexes := make([]*invertedindex.InvertedIndex, 0)

	for _, cond := range conds {
		iter := snapshot.NewIterator(nil, nil)
		for iter.Seek(cond.Start().Raw()); iter.Valid(); iter.Next() {
			key := iter.Key()
			termKey := getTerm(key)
			segmentID := getSegmentID(key)
			if (!cond.IncludeEnd() && term.Comparator(cond.End(), termKey) == 0) ||
				term.Comparator(cond.End(), termKey) < 0 {
				break
			}
			if !cond.IncludeStart() && term.Comparator(cond.Start(), termKey) == 0 {
				continue
			}

			for i := len(indexes); i <= int(segmentID); i++ {
				indexes = append(indexes, invertedindex.NewInvertedIndex(i))
			}
			index, err := fieldindex.Deserialize(iter.Value())
			if err != nil {
				return nil, fmt.Errorf("fail to get index iterator: %w", err)
			}
			indexes[segmentID].Put(termKey, index)
		}
		iter.Release()
		if err := iter.Error(); err != nil {
			return nil, fmt.Errorf("iter error: %w", err)
		}
	}

	indexIterator := NewIndexIterator(indexes)

	return indexIterator, nil
}

func (m *lsmStorage) Close() error {
	return m.index.Close()
}
