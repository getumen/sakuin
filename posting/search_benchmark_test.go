package posting

import (
	"math/rand"
	"sort"
	"testing"
)

/*
BenchmarkBinarySearchForPosting

goos: linux
goarch: amd64
pkg: github.com/getumen/sakuin/posting
cpu: Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
BenchmarkBinarySearch
BenchmarkBinarySearch/sort.Slice
BenchmarkBinarySearch/sort.Slice-4         	17703418	        57.09 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearch/standard
BenchmarkBinarySearch/standard-4           	67655511	        18.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearch/branch-prediction
BenchmarkBinarySearch/branch-prediction-4  	39635217	        27.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearch/branchless
BenchmarkBinarySearch/branchless-4         	16627096	        69.78 ns/op	       0 B/op	       0 allocs/op
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
BenchmarkBinarySearchForPosition/sort.Slice-4         	19735444	        58.75 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearchForPosition/standard
BenchmarkBinarySearchForPosition/standard-4
68894480	        16.70 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearchForPosition/branch-prediction
BenchmarkBinarySearchForPosition/branch-prediction-4  	34778972	        29.95 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinarySearchForPosition/branchless
BenchmarkBinarySearchForPosition/branchless-4         	22366386	        48.42 ns/op	       0 B/op	       0 allocs/op
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
