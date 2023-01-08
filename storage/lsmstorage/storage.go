package lsmstorage

import (
	"context"
	"fmt"

	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/invertedindex"
	"github.com/getumen/sakuin/storage"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
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

func (m *lsmStorage) Merge(
	ctx context.Context,
	index *invertedindex.InvertedIndex,
) error {
	tx, err := m.index.OpenTransaction()
	if err != nil {
		return fmt.Errorf("fail to open tx: %w", err)
	}
	it := index.Iterator()

	terms := make([]term.Term, 0)
	segments := make([]*Segment, 0)

	var maxSegmentID int
	for it.Next() {
		key := it.Key().(term.Term)
		terms = append(terms, key)

		value := it.Value().(fieldindex.FieldIndex)

		diskIndexBlob, err := tx.Get(key.Raw(), nil)
		if err == leveldb.ErrNotFound {
			seg := NewSegment(make([]byte, 0))
			segments = append(segments, seg)
			maxSegmentID = max(maxSegmentID, seg.FindAvailableSegment(value.EstimateSize(), maxSegmentSize))
		} else if err != nil {
			return fmt.Errorf("fail to get index: %w", err)
		} else {
			seg := NewSegment(diskIndexBlob)
			segments = append(segments, seg)
			maxSegmentID = max(maxSegmentID, seg.FindAvailableSegment(value.EstimateSize(), maxSegmentSize))
		}
	}

	it = index.Iterator()
	batch := new(leveldb.Batch)

	for i := 0; i < len(terms) && it.Next(); i++ {
		memoryFieldIndex := it.Value().(fieldindex.FieldIndex)
		fieldIndex, err := segments[i].Get(maxSegmentID)
		if err != nil {
			return fmt.Errorf("fail to merge index: %w", err)
		}
		fieldIndex.Merge(memoryFieldIndex)
		err = segments[i].Save(maxSegmentID, fieldIndex)
		if err != nil {
			return fmt.Errorf("fail to merge index: %w", err)
		}
		batch.Put(terms[i], segments[i].Bytes())
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

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func (m *lsmStorage) GetIndexIterator(
	ctx context.Context,
	conds []*termcond.TermCondition,
) (storage.IndexIterator, error) {
	snapshot, err := m.index.GetSnapshot()
	if err != nil {
		return nil, fmt.Errorf("fail to get snapshot: %w", err)
	}

	terms := make([]term.Term, 0)
	segmentIterators := make([]*SegmentIterator, 0)

	for _, cond := range conds {
		iter := snapshot.NewIterator(nil, nil)
		for iter.Seek(cond.Start().Raw()); iter.Valid(); iter.Next() {
			key := copyBytes(iter.Key())
			termKey := term.Term(key)
			if (!cond.IncludeEnd() && term.Comparator(cond.End(), termKey) == 0) ||
				term.Comparator(cond.End(), termKey) < 0 {
				break
			}
			if !cond.IncludeStart() && term.Comparator(cond.Start(), termKey) == 0 {
				continue
			}
			// cond contains key
			terms = append(terms, termKey)
			seg := NewSegment(copyBytes(iter.Value()))
			segmentIterators = append(segmentIterators, seg.Iterator())
		}
		iter.Release()
		if err := iter.Error(); err != nil {
			return nil, fmt.Errorf("iter error: %w", err)
		}
	}

	indexIterator := NewIndexIterator(terms, segmentIterators)

	return indexIterator, nil
}

func copyBytes(value []byte) []byte {
	result := make([]byte, len(value))
	copy(result, value)
	return result
}

func (m *lsmStorage) Close() error {
	return m.index.Close()
}
