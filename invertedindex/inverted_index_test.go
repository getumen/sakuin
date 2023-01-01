package invertedindex_test

import (
	"testing"

	"github.com/getumen/sakuin/expression"
	"github.com/getumen/sakuin/fieldindex"
	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/invertedindex"
	"github.com/getumen/sakuin/position"
	"github.com/getumen/sakuin/posting"
	"github.com/getumen/sakuin/postinglist"
	"github.com/getumen/sakuin/term"
	"github.com/getumen/sakuin/termcond"
	"github.com/stretchr/testify/require"
)

func TestInvertedIndex_Search(t *testing.T) {
	type args struct {
		exp *expression.Expression
	}
	tests := []struct {
		name string
		get  func() *invertedindex.InvertedIndex
		args args
		want *postinglist.PostingList
	}{
		{
			name: "search ab",
			get: func() *invertedindex.InvertedIndex {
				result := invertedindex.NewInvertedIndex()
				result.Put(
					term.NewText("a"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList(
								[]*posting.Posting{
									posting.NewPosting(1, position.NewPositions([]int64{1, 5, 7})),
									posting.NewPosting(2, position.NewPositions([]int64{3, 5, 7})),
								},
							),
						},
					),
				)
				result.Put(
					term.NewText("b"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList(
								[]*posting.Posting{
									posting.NewPosting(1, position.NewPositions([]int64{2})),
									posting.NewPosting(2, position.NewPositions([]int64{2})),
								},
							),
						},
					),
				)

				return result
			},
			args: args{
				exp: expression.NewPhrase(
					[]*expression.Expression{
						expression.NewFeature(
							expression.NewFeatureSpec(
								"f1",
								termcond.NewEqual(term.NewText("a"))),
						),
						expression.NewFeature(
							expression.NewFeatureSpec(
								"f1",
								termcond.NewEqual(term.NewText("b"))),
						),
					},
					[]int64{0, 1},
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(1, position.NewPositions([]int64{1})),
			}),
		},
		{
			name: "search a and not b",
			get: func() *invertedindex.InvertedIndex {
				result := invertedindex.NewInvertedIndex()
				result.Put(
					term.NewText("a"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList(
								[]*posting.Posting{
									posting.NewPosting(1, position.NewPositions([]int64{1, 5, 7})),
									posting.NewPosting(2, position.NewPositions([]int64{1, 5, 7})),
								},
							),
						},
					),
				)
				result.Put(
					term.NewText("b"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList(
								[]*posting.Posting{
									posting.NewPosting(1, position.NewPositions([]int64{2})),
								},
							),
						},
					),
				)

				return result
			},
			args: args{
				exp: expression.NewAnd(
					[]*expression.Expression{
						expression.NewFeature(
							expression.NewFeatureSpec(
								"f1",
								termcond.NewEqual(term.NewText("a"))),
						),
						expression.NewNot(
							expression.NewFeature(
								expression.NewFeatureSpec(
									"f1", termcond.NewEqual(term.NewText("b"))),
							),
						),
					},
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(2, position.NewPositions([]int64{1, 5, 7})),
			}),
		},
		{
			name: "search a or b",
			get: func() *invertedindex.InvertedIndex {
				result := invertedindex.NewInvertedIndex()
				result.Put(
					term.NewText("a"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{1, 5, 7})),
								posting.NewPosting(2, position.NewPositions([]int64{1, 5, 7})),
							}),
						},
					),
				)
				result.Put(
					term.NewText("b"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(3, position.NewPositions([]int64{2})),
							}),
						},
					),
				)

				return result
			},
			args: args{
				exp: expression.NewOr(
					[]*expression.Expression{
						expression.NewFeature(
							expression.NewFeatureSpec(
								"f1",
								termcond.NewEqual(term.NewText("a"))),
						),
						expression.NewFeature(
							expression.NewFeatureSpec(
								"f1",
								termcond.NewEqual(term.NewText("b"))),
						),
					},
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(1, position.NewPositions([]int64{1, 5, 7})),
				posting.NewPosting(2, position.NewPositions([]int64{1, 5, 7})),
				posting.NewPosting(3, position.NewPositions([]int64{2})),
			}),
		},
		{
			name: "search ab and not cd",
			get: func() *invertedindex.InvertedIndex {
				result := invertedindex.NewInvertedIndex()
				result.Put(
					term.NewText("a"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{1, 5, 7})),
								posting.NewPosting(2, position.NewPositions([]int64{1, 5, 7})),
							}),
						},
					),
				)
				result.Put(
					term.NewText("b"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{2})),
								posting.NewPosting(2, position.NewPositions([]int64{2})),
							}),
						},
					),
				)
				result.Put(
					term.NewText("c"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{2})),
							}),
						},
					),
				)
				result.Put(
					term.NewText("d"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{3})),
							}),
						},
					),
				)

				return result
			},
			args: args{
				exp: expression.NewAnd(
					[]*expression.Expression{
						expression.NewPhrase([]*expression.Expression{
							expression.NewFeature(
								expression.NewFeatureSpec(
									"f1",
									termcond.NewEqual(term.NewText("a"))),
							),
							expression.NewFeature(
								expression.NewFeatureSpec(
									"f1",
									termcond.NewEqual(term.NewText("b"))),
							),
						},
							[]int64{0, 1},
						),
						expression.NewNot(
							expression.NewPhrase([]*expression.Expression{
								expression.NewFeature(
									expression.NewFeatureSpec(
										"f1",
										termcond.NewEqual(term.NewText("c"))),
								),
								expression.NewFeature(
									expression.NewFeatureSpec(
										"f1",
										termcond.NewEqual(term.NewText("d"))),
								),
							},
								[]int64{0, 1},
							),
						),
					},
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(2, position.NewPositions([]int64{1})),
			}),
		},
		{
			name: "search (ab and not cd) or e",
			get: func() *invertedindex.InvertedIndex {
				result := invertedindex.NewInvertedIndex()
				result.Put(
					term.NewText("a"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{1, 5, 7})),
								posting.NewPosting(2, position.NewPositions([]int64{1, 5, 7})),
							}),
						},
					),
				)
				result.Put(
					term.NewText("b"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{2})),
								posting.NewPosting(2, position.NewPositions([]int64{2})),
							}),
						},
					),
				)
				result.Put(
					term.NewText("c"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{2})),
							}),
						},
					),
				)
				result.Put(
					term.NewText("d"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{3})),
							}),
						},
					),
				)
				result.Put(
					term.NewText("e"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(3, position.NewPositions([]int64{1})),
							}),
						},
					),
				)

				return result
			},
			args: args{
				exp: expression.NewOr(
					[]*expression.Expression{
						expression.NewAnd(
							[]*expression.Expression{
								expression.NewPhrase([]*expression.Expression{
									expression.NewFeature(
										expression.NewFeatureSpec(
											"f1",
											termcond.NewEqual(term.NewText("a"))),
									),
									expression.NewFeature(
										expression.NewFeatureSpec(
											"f1",
											termcond.NewEqual(term.NewText("b"))),
									),
								},
									[]int64{0, 1},
								),
								expression.NewNot(
									expression.NewPhrase([]*expression.Expression{
										expression.NewFeature(
											expression.NewFeatureSpec(
												"f1",
												termcond.NewEqual(term.NewText("c"))),
										),
										expression.NewFeature(
											expression.NewFeatureSpec(
												"f1",
												termcond.NewEqual(term.NewText("d"))),
										),
									},
										[]int64{0, 1},
									),
								),
							},
						),
						expression.NewFeature(
							expression.NewFeatureSpec(
								"f1",
								termcond.NewEqual(term.NewText("e")))),
					},
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(2, position.NewPositions([]int64{1})),
				posting.NewPosting(3, position.NewPositions([]int64{1})),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.get()
			got := target.Search(tt.args.exp)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestInvertedIndex_GetPostingListInFeature(t *testing.T) {
	type args struct {
		feature *expression.FeatureSpec
	}
	tests := []struct {
		name string
		get  func() *invertedindex.InvertedIndex
		args args
		want *postinglist.PostingList
	}{
		{
			name: "equal",
			get: func() *invertedindex.InvertedIndex {
				result := invertedindex.NewInvertedIndex()
				result.Put(
					term.NewText("a"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{1, 5, 7})),
								posting.NewPosting(2, position.NewPositions([]int64{1, 5, 7})),
							}),
						},
					),
				)
				result.Put(
					term.NewText("b"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{2})),
								posting.NewPosting(2, position.NewPositions([]int64{2})),
							}),
						},
					),
				)

				return result
			},
			args: args{
				feature: expression.NewFeatureSpec(
					"f1",
					termcond.NewEqual(term.NewText("a")),
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(1, position.NewPositions([]int64{1, 5, 7})),
				posting.NewPosting(2, position.NewPositions([]int64{1, 5, 7})),
			}),
		},
		{
			name: "range",
			get: func() *invertedindex.InvertedIndex {
				result := invertedindex.NewInvertedIndex()
				result.Put(
					term.NewText("a"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{1, 5, 7})),
								posting.NewPosting(2, position.NewPositions([]int64{1, 5, 7})),
							}),
						},
					),
				)
				result.Put(
					term.NewText("b"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(1, position.NewPositions([]int64{2})),
								posting.NewPosting(3, position.NewPositions([]int64{2})),
							}),
						},
					),
				)
				result.Put(
					term.NewText("A"),
					fieldindex.NewFieldIndexFromMap(
						map[fieldname.FieldName]*postinglist.PostingList{
							"f1": postinglist.NewPostingList([]*posting.Posting{
								posting.NewPosting(4, position.NewPositions([]int64{2})),
							}),
						},
					),
				)

				return result
			},
			args: args{
				feature: expression.NewFeatureSpec(
					"f1",
					termcond.NewRange(
						term.NewText("a"),
						true,
						term.NewText("z"),
						true,
					),
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(1, position.NewPositions([]int64{1, 2, 5, 7})),
				posting.NewPosting(2, position.NewPositions([]int64{1, 5, 7})),
				posting.NewPosting(3, position.NewPositions([]int64{2})),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := tt.get()
			got := i.GetPostingListInFeature(tt.args.feature)
			require.Equal(t, tt.want, got)
		})
	}
}
