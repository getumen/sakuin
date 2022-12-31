package term

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func FuzzInt64Bytes_Order(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y int64) {
		xBytes := newInt64Bytes(x)
		yBytes := newInt64Bytes(y)
		require.Equal(t, x > y, xBytes > yBytes)
		require.Equal(t, x == y, xBytes == yBytes)
	})
}

func FuzzInt64Bytes_Bytes(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y int64) {
		xBytes := newInt64Bytes(x)
		yBytes := newInt64Bytes(y)
		require.Equal(t, x > y, bytes.Compare(xBytes.bytes(), yBytes.bytes()) > 0)
		require.Equal(t, x < y, bytes.Compare(xBytes.bytes(), yBytes.bytes()) < 0)
	})
}

func FuzzInt64Bytes_Float64(f *testing.F) {
	f.Fuzz(func(t *testing.T, x int64) {
		require.Equal(t, newInt64Bytes(x).int64(), x)
	})
}

func FuzzFloat64Bytes_Order(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y float64) {
		xBytes := newFloat64Bytes(x)
		yBytes := newFloat64Bytes(y)
		require.Equal(t, x > y, xBytes > yBytes)
		require.Equal(t, x == y, xBytes == yBytes)
	})
}

func FuzzFloat64Bytes_Bytes(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y float64) {
		xBytes := newFloat64Bytes(x)
		yBytes := newFloat64Bytes(y)
		require.Equal(t, x > y, bytes.Compare(xBytes.bytes(), yBytes.bytes()) > 0)
		require.Equal(t, x < y, bytes.Compare(xBytes.bytes(), yBytes.bytes()) < 0)
	})
}

func FuzzFloat64Bytes_Float64(f *testing.F) {
	f.Fuzz(func(t *testing.T, x float64) {
		require.Equal(t, newFloat64Bytes(x).float64(), x)
	})
}

func FuzzSignFlip(f *testing.F) {
	f.Fuzz(func(t *testing.T, x float64) {
		require.Equal(t, math.Float64frombits(signFlip(math.Float64bits(x))), -x)
	})
}
