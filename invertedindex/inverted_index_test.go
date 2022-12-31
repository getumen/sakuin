package invertedindex_test

import (
	"testing"

	"github.com/getumen/sakuin/booleanexpression"
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
		booleanExpression *booleanexpression.BooleanExpression
	}
	tests := []struct {
		name string
		get  func() *invertedindex.InvertedIndex
		args args
		want *postinglist.PostingList
	}{
		{
			name: "search a and b",
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
									posting.NewPosting(2, position.NewPositions([]int64{2})),
								},
							),
						},
					),
				)

				return result
			},
			args: args{
				booleanExpression: booleanexpression.NewAnd(
					[]*booleanexpression.BooleanExpression{
						booleanexpression.NewFeature(
							booleanexpression.NewBoolenaFeature(
								"f1",
								termcond.NewEqual(term.NewText("a"))),
						),
						booleanexpression.NewFeature(
							booleanexpression.NewBoolenaFeature("f1", termcond.NewEqual(term.NewText("b"))),
						),
					},
					[]int64{0, 1},
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(1, nil),
				posting.NewPosting(2, nil),
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
				booleanExpression: booleanexpression.NewAnd(
					[]*booleanexpression.BooleanExpression{
						booleanexpression.NewFeature(
							booleanexpression.NewBoolenaFeature(
								"f1",
								termcond.NewEqual(term.NewText("a"))),
						),
						booleanexpression.NewNot(
							booleanexpression.NewFeature(
								booleanexpression.NewBoolenaFeature(
									"f1", termcond.NewEqual(term.NewText("b"))),
							),
						),
					},
					nil,
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(2, nil),
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
				booleanExpression: booleanexpression.NewOr(
					[]*booleanexpression.BooleanExpression{
						booleanexpression.NewFeature(
							booleanexpression.NewBoolenaFeature(
								"f1",
								termcond.NewEqual(term.NewText("a"))),
						),
						booleanexpression.NewFeature(
							booleanexpression.NewBoolenaFeature(
								"f1",
								termcond.NewEqual(term.NewText("b"))),
						),
					},
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(1, nil),
				posting.NewPosting(2, nil),
				posting.NewPosting(3, nil),
			}),
		},
		{
			name: "search (a and b) and not (c and d)",
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
				booleanExpression: booleanexpression.NewAnd(
					[]*booleanexpression.BooleanExpression{
						booleanexpression.NewAnd([]*booleanexpression.BooleanExpression{
							booleanexpression.NewFeature(
								booleanexpression.NewBoolenaFeature(
									"f1",
									termcond.NewEqual(term.NewText("a"))),
							),
							booleanexpression.NewFeature(
								booleanexpression.NewBoolenaFeature(
									"f1",
									termcond.NewEqual(term.NewText("b"))),
							),
						},
							[]int64{0, 1},
						),
						booleanexpression.NewNot(
							booleanexpression.NewAnd([]*booleanexpression.BooleanExpression{
								booleanexpression.NewFeature(
									booleanexpression.NewBoolenaFeature(
										"f1",
										termcond.NewEqual(term.NewText("c"))),
								),
								booleanexpression.NewFeature(
									booleanexpression.NewBoolenaFeature(
										"f1",
										termcond.NewEqual(term.NewText("d"))),
								),
							},
								[]int64{0, 1},
							),
						),
					},
					nil,
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(2, nil),
			}),
		},
		{
			name: "search ((a and b) and not (c and d)) or e",
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
				booleanExpression: booleanexpression.NewOr(
					[]*booleanexpression.BooleanExpression{
						booleanexpression.NewAnd(
							[]*booleanexpression.BooleanExpression{
								booleanexpression.NewAnd([]*booleanexpression.BooleanExpression{
									booleanexpression.NewFeature(
										booleanexpression.NewBoolenaFeature(
											"f1",
											termcond.NewEqual(term.NewText("a"))),
									),
									booleanexpression.NewFeature(
										booleanexpression.NewBoolenaFeature(
											"f1",
											termcond.NewEqual(term.NewText("b"))),
									),
								},
									[]int64{0, 1},
								),
								booleanexpression.NewNot(
									booleanexpression.NewAnd([]*booleanexpression.BooleanExpression{
										booleanexpression.NewFeature(
											booleanexpression.NewBoolenaFeature(
												"f1",
												termcond.NewEqual(term.NewText("c"))),
										),
										booleanexpression.NewFeature(
											booleanexpression.NewBoolenaFeature(
												"f1",
												termcond.NewEqual(term.NewText("d"))),
										),
									},
										[]int64{0, 1},
									),
								),
							},
							nil,
						),
						booleanexpression.NewFeature(
							booleanexpression.NewBoolenaFeature(
								"f1",
								termcond.NewEqual(term.NewText("e")))),
					},
				),
			},
			want: postinglist.NewPostingList([]*posting.Posting{
				posting.NewPosting(2, nil),
				posting.NewPosting(3, nil),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.get()
			got := target.Search(tt.args.booleanExpression)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestInvertedIndex_GetPostingListInFeature(t *testing.T) {
	type args struct {
		feature *booleanexpression.BooleanFeature
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
				feature: booleanexpression.NewBoolenaFeature(
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
				feature: booleanexpression.NewBoolenaFeature(
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
				posting.NewPosting(1, nil),
				posting.NewPosting(2, nil),
				posting.NewPosting(3, nil),
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
