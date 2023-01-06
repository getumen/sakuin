package invertedindex

import (
	"reflect"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/getumen/sakuin/expression"
	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
)

type InvertedIndex struct {
	// key is term.Term
	// value is fieldindex.FieldIndex
	*redblacktree.Tree
	segmentID int
}

func NewInvertedIndex(segmentID int) *InvertedIndex {
	return &InvertedIndex{
		Tree:      redblacktree.NewWith(term.Comparator),
		segmentID: segmentID,
	}
}

func (i InvertedIndex) SegmentID() int {
	return i.segmentID
}

func (i *InvertedIndex) Equal(other *InvertedIndex) bool {

	if i.segmentID != other.segmentID {
		return false
	}
	if i.Size() != other.Size() {
		return false
	}
	xIter := i.Iterator()
	yIter := other.Iterator()
	for xIter.Next() && yIter.Next() {
		if i.Comparator(xIter.Key(), yIter.Key()) != 0 {

			return false
		}
		if !reflect.DeepEqual(xIter.Value(), yIter.Value()) {
			return false
		}
	}
	return true
}

func (i *InvertedIndex) Merge(other *InvertedIndex) {
	otherIter := other.Iterator()
	for otherIter.Next() {
		otherKey := otherIter.Key().(term.Term)
		otherValue := otherIter.Value().(fieldindex.FieldIndex)
		if fieldIndex, ok := i.Get(otherKey); !ok {
			fieldIndex = otherValue
			i.Put(otherKey, fieldIndex)
		} else {
			fieldIndex.(fieldindex.FieldIndex).Merge(otherValue)
		}
	}
}

func (i InvertedIndex) GetPostingListInFeature(feature *expression.FeatureSpec) *postinglist.PostingList {
	field := feature.Field()
	cond := feature.TermCondition()
	node, ok := i.Ceiling(cond.Start())
	if !ok {
		return postinglist.NewPostingList(make([]*posting.Posting, 0))
	}
	it := i.IteratorAt(node)
	postingLists := make([]*postinglist.PostingList, 0)

	for {
		key := it.Key().(term.Term)
		if !cond.IncludeStart() && term.Comparator(key, cond.Start()) == 0 {
			continue
		}
		if (term.Comparator(key, cond.End()) == 0 && !cond.IncludeEnd()) ||
			term.Comparator(key, cond.End()) > 0 {
			break
		}
		fieldIndex := it.Value().(fieldindex.FieldIndex)
		if postingList, ok := fieldIndex[field]; ok {
			postingLists = append(postingLists, postingList)
		}
		if !it.Next() {
			break
		}
	}

	return postinglist.Union(postingLists)
}

func (i InvertedIndex) Search(exp *expression.Expression) *postinglist.PostingList {
	if exp.And() != nil {
		sets := make([]*postinglist.PostingList, 0)
		excludeSets := make([]*postinglist.PostingList, 0)
		for _, be := range exp.And() {
			if be.Not() != nil {
				set := i.Search(be.Not())
				excludeSets = append(excludeSets, set)
			} else {
				set := i.Search(be)
				sets = append(sets, set)
			}
		}

		return postinglist.Difference(
			postinglist.Intersection(sets),
			postinglist.Intersection(excludeSets),
		)

	}

	if exp.Phrase() != nil {
		sets := make([]*postinglist.PostingList, 0)
		for _, be := range exp.Phrase() {
			set := i.Search(be)
			sets = append(sets, set)
		}

		return postinglist.PhraseMatch(sets, exp.RelativePosition())
	}

	if exp.Or() != nil {
		sets := make([]*postinglist.PostingList, 0)
		for _, be := range exp.Or() {
			if be.Not() != nil {
				continue
			}
			set := i.Search(be)
			sets = append(sets, set)
		}
		return postinglist.Union(sets)
	}
	if exp.Feature() != nil {
		feature := exp.Feature()
		postingList := i.GetPostingListInFeature(feature)
		return postingList
	}
	return postinglist.NewPostingList(make([]*posting.Posting, 0))
}

func (i InvertedIndex) GetPartialIndex(conds []*termcond.TermCondition) *InvertedIndex {
	result := NewInvertedIndex(i.segmentID)
	for _, cond := range conds {
		node, ok := i.Ceiling(cond.Start())
		if !ok {
			continue
		}
		it := i.IteratorAt(node)

		for {
			key := it.Key().(term.Term)
			if !cond.IncludeStart() && term.Comparator(key, cond.Start()) == 0 {
				continue
			}
			if (term.Comparator(key, cond.End()) == 0 && !cond.IncludeEnd()) ||
				term.Comparator(key, cond.End()) > 0 {
				break
			}
			fieldIndex := it.Value().(fieldindex.FieldIndex)
			result.Put(key, fieldIndex)
			if !it.Next() {
				break
			}
		}
	}
	return result
}
