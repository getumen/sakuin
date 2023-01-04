package writer

import (
	"context"
	"fmt"

	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/invertedindex"
	"github.com/getumen/sakuin/storage"
)

type IndexWriter struct {
	storage storage.IndexStorage
}

func NewIndexWriter(
	storage storage.IndexStorage,
) *IndexWriter {
	return &IndexWriter{
		storage: storage,
	}
}

func (s *IndexWriter) CreateDocuments(
	ctx context.Context,
	docs []*document.Document,
) error {
	builder := invertedindex.NewBuilder()
	for _, doc := range docs {
		builder.AddDocument(doc)
	}
	index := builder.Build()

	if err := s.commit(ctx, index); err != nil {
		return err
	}
	return nil
}

func (s *IndexWriter) commit(
	ctx context.Context,
	index *invertedindex.InvertedIndex,
) error {
	err := s.storage.Merge(ctx, index)
	if err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	return nil
}
