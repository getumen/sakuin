package posting

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinarySearchForPosting(t *testing.T) {
	type args struct {
		sortedArray []*Posting
		posting     *Posting
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "found",
			args: args{
				sortedArray: []*Posting{
					NewPosting(1, nil),
					NewPosting(3, nil),
					NewPosting(6, nil),
					NewPosting(10, nil),
					NewPosting(15, nil),
					NewPosting(21, nil),
					NewPosting(28, nil),
					NewPosting(36, nil),
					NewPosting(45, nil),
					NewPosting(55, nil),
				},
				posting: NewPosting(6, nil),
			},
			want: 2,
		},
		{
			name: "not found",
			args: args{
				sortedArray: []*Posting{
					NewPosting(1, nil),
					NewPosting(3, nil),
					NewPosting(6, nil),
					NewPosting(10, nil),
					NewPosting(15, nil),
					NewPosting(21, nil),
					NewPosting(28, nil),
					NewPosting(36, nil),
					NewPosting(45, nil),
					NewPosting(55, nil),
				},
				posting: NewPosting(14, nil),
			},
			want: 4,
		},
	}
	for _, algirthm := range []struct {
		name string
		alg  func([]*Posting, *Posting) int
	}{
		{
			name: "sort.Slice",
			alg:  sortBinarySearchForPosting,
		},
		{
			name: "standard",
			alg:  standardBinarySearchForPosting,
		},
		{
			name: "branch-prediction",
			alg:  branchPredictionBinarySearchForPosting,
		},
		{
			name: "branchless",
			alg:  branchlessBinarySearchForPosting,
		},
	} {
		for _, tt := range tests {
			t.Run(algirthm.name+" "+tt.name, func(t *testing.T) {
				got := algirthm.alg(tt.args.sortedArray, tt.args.posting)
				require.Equal(t, tt.want, got)
			})
		}
	}
}

func TestBinarySearchForPosition(t *testing.T) {
	type args struct {
		sortedArray []uint32
		posting     uint32
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "found",
			args: args{
				sortedArray: []uint32{
					1,
					3,
					6,
					10,
					15,
					21,
					28,
					36,
					45,
					55,
				},
				posting: 6,
			},
			want: 2,
		},
		{
			name: "not found",
			args: args{
				sortedArray: []uint32{
					1,
					3,
					6,
					10,
					15,
					21,
					28,
					36,
					45,
					55,
				},
				posting: 14,
			},
			want: 4,
		},
	}
	for _, algirthm := range []struct {
		name string
		alg  func([]uint32, uint32) int
	}{
		{
			name: "sort.Slice",
			alg:  sortBinarySearchForPosition,
		},
		{
			name: "standard",
			alg:  standardBinarySearchForPosition,
		},
		{
			name: "branch-prediction",
			alg:  branchPredictionBinarySearchForPosition,
		},
		{
			name: "branchless",
			alg:  branchlessBinarySearchForPosition,
		},
	} {
		for _, tt := range tests {
			t.Run(algirthm.name+" "+tt.name, func(t *testing.T) {
				got := algirthm.alg(tt.args.sortedArray, tt.args.posting)
				require.Equal(t, tt.want, got)
			})
		}
	}
}

func TestExponentialSearch(t *testing.T) {
	type args struct {
		sortedArray []*Posting
		posting     *Posting
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "found",
			args: args{
				sortedArray: []*Posting{
					NewPosting(1, nil),
					NewPosting(3, nil),
					NewPosting(6, nil),
					NewPosting(10, nil),
					NewPosting(15, nil),
					NewPosting(21, nil),
					NewPosting(28, nil),
					NewPosting(36, nil),
					NewPosting(45, nil),
					NewPosting(55, nil),
				},
				posting: NewPosting(6, nil),
			},
			want: 2,
		},
		{
			name: "not found",
			args: args{
				sortedArray: []*Posting{
					NewPosting(1, nil),
					NewPosting(3, nil),
					NewPosting(6, nil),
					NewPosting(10, nil),
					NewPosting(15, nil),
					NewPosting(21, nil),
					NewPosting(28, nil),
					NewPosting(36, nil),
					NewPosting(45, nil),
					NewPosting(55, nil),
				},
				posting: NewPosting(14, nil),
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExponentialSearch(tt.args.sortedArray, tt.args.posting)
			require.Equal(t, tt.want, got)
		})
	}
}

/*
BenchmarkBinarySearchForPosting

goos: linux
goarch: amd64
pkg: github.com/getumen/sakuin/posting
cpu: Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
BenchmarkBinarySearchForPosting
BenchmarkBinarySearchForPosting/sort.Slice
BenchmarkBinarySearchForPosting/sort.Slice-4         	25396570	        48.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearchForPosting/standard
BenchmarkBinarySearchForPosting/standard-4           	62686153	        19.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearchForPosting/branch-prediction
BenchmarkBinarySearchForPosting/branch-prediction-4  	44398458	        24.13 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearchForPosting/branchless
BenchmarkBinarySearchForPosting/branchless-4         	15233991	        69.40 ns/op	       0 B/op	       0 allocs/op
PASS
*/
func BenchmarkBinarySearchForPosting(b *testing.B) {
	const dataNum = 100000
	data := make([]uint64, dataNum)
	for i := 0; i < dataNum; i++ {
		data[i] = rand.Uint64()
	}
	sort.Slice(data, func(i, j int) bool { return data[i] < data[j] })

	postingList := make([]*Posting, dataNum)
	for i, v := range data {
		postingList[i] = NewPosting(v, make([]uint32, 0))
	}

	for _, c := range []struct {
		name string
		alg  func(arr []*Posting, posting *Posting) int
	}{
		{
			name: "sort.Slice",
			alg:  sortBinarySearchForPosting,
		},
		{
			name: "standard",
			alg:  standardBinarySearchForPosting,
		},
		{
			name: "branch-prediction",
			alg:  branchPredictionBinarySearchForPosting,
		},
		{
			name: "branchless",
			alg:  branchlessBinarySearchForPosting,
		},
	} {
		b.Run(c.name, func(b *testing.B) {
			v := NewPosting(1000, make([]uint32, 0))
			for i := 0; i < b.N; i++ {
				c.alg(postingList, v)
			}
		})
	}

}

/*
BenchmarkBinarySearchForPosition

goos: linux
goarch: amd64
pkg: github.com/getumen/sakuin/posting
cpu: Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
BenchmarkBinarySearchForPosition
BenchmarkBinarySearchForPosition/sort.Slice
BenchmarkBinarySearchForPosition/sort.Slice-4         	26721198	        42.52 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearchForPosition/standard
BenchmarkBinarySearchForPosition/standard-4           	69359968	        17.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearchForPosition/branch-prediction
BenchmarkBinarySearchForPosition/branch-prediction-4  	47127693	        22.64 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearchForPosition/branchless
BenchmarkBinarySearchForPosition/branchless-4         	23144769	        46.94 ns/op	       0 B/op	       0 allocs/op
PASS
*/
func BenchmarkBinarySearchForPosition(b *testing.B) {
	const dataNum = 100000
	data := make([]uint32, dataNum)
	for i := 0; i < dataNum; i++ {
		data[i] = rand.Uint32()
	}
	sort.Slice(data, func(i, j int) bool { return data[i] < data[j] })

	for _, c := range []struct {
		name string
		alg  func(arr []uint32, posting uint32) int
	}{
		{
			name: "sort.Slice",
			alg:  sortBinarySearchForPosition,
		},
		{
			name: "standard",
			alg:  standardBinarySearchForPosition,
		},
		{
			name: "branch-prediction",
			alg:  branchPredictionBinarySearchForPosition,
		},
		{
			name: "branchless",
			alg:  branchlessBinarySearchForPosition,
		},
	} {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.alg(data, 1000)
			}
		})
	}

}
