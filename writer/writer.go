package writer

import (
	"context"
	"fmt"

	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/invertedindex"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/storage"
)

type IndexWriter struct {
	index   *invertedindex.InvertedIndex
	storage storage.IndexStorage
}

func NewIndexWriter(
	storage storage.IndexStorage,
) *IndexWriter {
	return &IndexWriter{
		index:   invertedindex.NewInvertedIndex(),
		storage: storage,
	}
}

func (s *IndexWriter) CreateDocuments(
	ctx context.Context,
	docs []*document.Document,
) error {
	for _, doc := range docs {
		s.addDocument(doc)
	}

	if err := s.commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *IndexWriter) addDocument(doc *document.Document) {
	index := s.indexDocument(doc)
	s.index.Merge(index)
}

func (s *IndexWriter) commit(ctx context.Context) error {
	err := s.storage.Merge(ctx, s.index)
	if err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	s.index = invertedindex.NewInvertedIndex()
	return nil
}

func (s *IndexWriter) indexDocument(
	doc *document.Document,
) *invertedindex.InvertedIndex {
	indexChunk := invertedindex.NewInvertedIndex()

	for _, field := range doc.Fields() {
		var pos uint32
		for _, term := range field.Content() {
			var fieldIndex fieldindex.FieldIndex

			if iFieldIndex, ok := indexChunk.Get(term); ok {
				fieldIndex = iFieldIndex.(fieldindex.FieldIndex)
			} else {
				fieldIndex = fieldindex.NewFieldIndex()
			}

			newPostingList := postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(doc.ID(), []uint32{pos}),
			})
			pos += 1

			postingList, ok := fieldIndex[field.FieldName()]
			if ok {
				newPostingList.Merge(postingList)
			}

			fieldIndex[field.FieldName()] = newPostingList

			indexChunk.Put(term, fieldIndex)
		}
	}

	return indexChunk
}
