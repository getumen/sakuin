package int64enc

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

const length = 10000

func TestSimple8B(t *testing.T) {
	const length = 10000
	expected := make([]uint64, length)
	for i := 0; i < length; i++ {
		expected[i] = rand.Uint64() % 256
	}
	blob, err := encodeSimple8B(expected)
	require.NoError(t, err)
	t.Logf("byte size: %d", len(blob))
	actual, err := decodeSimple8B(blob)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func BenchmarkEncodeSimple8B(b *testing.B) {
	rand.Seed(42)
	expected := make([]uint64, length)
	for i := 0; i < length; i++ {
		expected[i] = rand.Uint64() % 256
	}

	b.Run("encode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = encodeSimple8B(expected)
		}
	})
}

func BenchmarkDecodeSimple8B(b *testing.B) {
	rand.Seed(42)
	expected := make([]uint64, length)
	for i := 0; i < length; i++ {
		expected[i] = rand.Uint64() % 256
	}
	blob, _ := encodeSimple8B(expected)

	b.Run("decode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = decodeSimple8B(blob)
		}
	})
}
