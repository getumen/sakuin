package int32enc

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

const length = 10000

func TestStreamVByte(t *testing.T) {
	const length = 10000
	expected := make([]uint32, length)
	for i := 0; i < length; i++ {
		expected[i] = rand.Uint32() % 256
	}
	blob, err := encodeSteamVbyte(expected)
	require.NoError(t, err)
	t.Logf("byte size: %d", len(blob))
	actual, err := decodeSteamVbyte(blob)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func BenchmarkEncodeStreamVByte(b *testing.B) {
	rand.Seed(42)
	expected := make([]uint32, length)
	for i := 0; i < length; i++ {
		expected[i] = rand.Uint32() % 256
	}

	b.Run("encode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = encodeSteamVbyte(expected)
		}
	})
}

func BenchmarkDecodeStreamVByte(b *testing.B) {
	rand.Seed(42)
	expected := make([]uint32, length)
	for i := 0; i < length; i++ {
		expected[i] = rand.Uint32() % 256
	}
	blob, _ := encodeSteamVbyte(expected)

	b.Run("decode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = decodeSteamVbyte(blob)
		}
	})
}
