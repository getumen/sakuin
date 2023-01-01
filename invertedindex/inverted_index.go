package invertedindex

import (
	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/getumen/sakuin/booleanexpression"
	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/term"
)

type InvertedIndex struct {
	// key is term.Term
	// value is fieldindex.FieldIndex
	*redblacktree.Tree
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Tree: redblacktree.NewWith(term.Comparator),
	}
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

func (i InvertedIndex) GetPostingListInFeature(feature *booleanexpression.BooleanFeature) *postinglist.PostingList {
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

func (i InvertedIndex) Search(booleanExpression *booleanexpression.BooleanExpression) *postinglist.PostingList {
	if booleanExpression.And() != nil {
		sets := make([]*postinglist.PostingList, 0)
		excludeSets := make([]*postinglist.PostingList, 0)
		for _, be := range booleanExpression.And() {
			if be.Not() != nil {
				set := i.Search(be.Not())
				excludeSets = append(excludeSets, set)
			} else {
				set := i.Search(be)
				sets = append(sets, set)
			}
		}

		excludeIntersection := postinglist.Intersection(excludeSets)

		if booleanExpression.RelativePosition() == nil {
			return postinglist.Difference(
				postinglist.Intersection(sets),
				excludeIntersection,
			)
		}

		return postinglist.Difference(
			postinglist.PhraseMatch(sets, booleanExpression.RelativePosition()),
			excludeIntersection,
		)

	}
	if booleanExpression.Or() != nil {
		sets := make([]*postinglist.PostingList, 0)
		for _, be := range booleanExpression.Or() {
			if be.Not() != nil {
				continue
			}
			set := i.Search(be)
			sets = append(sets, set)
		}
		return postinglist.Union(sets)
	}
	if booleanExpression.Feature() != nil {
		feature := booleanExpression.Feature()
		postingList := i.GetPostingListInFeature(feature)
		return postingList
	}
	return postinglist.NewPostingList(make([]*posting.Posting, 0))
}
