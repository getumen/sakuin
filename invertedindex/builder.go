package invertedindex

import (
	"fmt"
	"sort"

	"github.com/getumen/sakuin/document"
	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/term"
)

type InvertedIndexBuilder struct {
	elements []*indexElement
}

func NewBuilder() *InvertedIndexBuilder {
	return &InvertedIndexBuilder{
		elements: make([]*indexElement, 0),
	}
}

func (b *InvertedIndexBuilder) AddDocument(doc *document.Document) {
	for _, f := range doc.Fields() {
		var pos uint32
		for _, t := range f.Content() {
			pos += 1
			b.elements = append(b.elements, NewIndexElement(
				f.FieldName(), t, doc.ID(), pos,
			))
		}
	}
}

func (b *InvertedIndexBuilder) Build() *InvertedIndex {
	sort.Slice(b.elements, func(i, j int) bool {
		return b.elements[i].compare(b.elements[j]) < 0
	})

	index := NewInvertedIndex()
	var lastIndexElement *indexElement
	fieldIndex := fieldindex.NewFieldIndex()
	postings := make([]*posting.Posting, 0)
	positions := make([]uint32, 0)

	for i := range b.elements {
		if lastIndexElement == nil {
			lastIndexElement = b.elements[i]
			continue
		}

		if term.Comparator(b.elements[i].term, lastIndexElement.term) != 0 {
			positions = append(positions, lastIndexElement.position)
			postings = append(postings, posting.NewPosting(
				lastIndexElement.docID, positions,
			))
			fieldIndex[lastIndexElement.field] = postinglist.NewPostingList(postings)
			index.Put(lastIndexElement.term, fieldIndex)
			lastIndexElement = b.elements[i]
			fieldIndex = fieldindex.NewFieldIndex()
			postings = make([]*posting.Posting, 0)
			positions = make([]uint32, 0)
			continue
		}

		if b.elements[i].field != lastIndexElement.field {
			positions = append(positions, lastIndexElement.position)
			postings = append(postings, posting.NewPosting(
				lastIndexElement.docID, positions,
			))
			fieldIndex[lastIndexElement.field] = postinglist.NewPostingList(postings)
			lastIndexElement = b.elements[i]
			postings = make([]*posting.Posting, 0)
			positions = make([]uint32, 0)
			continue
		}

		if b.elements[i].docID != lastIndexElement.docID {
			positions = append(positions, lastIndexElement.position)
			postings = append(postings, posting.NewPosting(
				lastIndexElement.docID, positions,
			))
			lastIndexElement = b.elements[i]
			positions = make([]uint32, 0)
			continue
		}

		if b.elements[i].position != lastIndexElement.position {
			positions = append(positions, lastIndexElement.position)
			lastIndexElement = b.elements[i]
			continue
		}

	}
	lastIndexElement = b.elements[len(b.elements)-1]
	positions = append(positions, lastIndexElement.position)
	postings = append(postings, posting.NewPosting(
		lastIndexElement.docID, positions,
	))
	fieldIndex[lastIndexElement.field] = postinglist.NewPostingList(postings)
	index.Put(lastIndexElement.term, fieldIndex)

	return index
}

type indexElement struct {
	field    fieldname.FieldName
	term     term.Term
	docID    uint64
	position uint32
}

func NewIndexElement(
	field fieldname.FieldName,
	term term.Term,
	docID uint64,
	position uint32,
) *indexElement {
	return &indexElement{
		field:    field,
		term:     term,
		docID:    docID,
		position: position,
	}
}

func (i *indexElement) String() string {
	return fmt.Sprintf(
		"{term:%s field:%s docID:%d position:%d}",
		i.field, i.term.String(), i.docID, i.position)
}

func (e *indexElement) compare(other *indexElement) int {
	if term.Comparator(e.term, other.term) != 0 {
		return term.Comparator(e.term, other.term)
	} else if e.field != other.field {
		if e.field < other.field {
			return -1
		} else if e.field == other.field {
			return 0
		} else {
			return 1
		}
	} else if e.docID != other.docID {
		return int(e.docID) - int(other.docID)
	} else {
		return int(e.position) - int(other.position)
	}
}
