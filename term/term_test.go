package term

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt64Bytes_Order(t *testing.T) {

	for _, c := range []struct {
		name string
		x    int64
		y    int64
	}{
		{
			name: "pos pos",
			x:    1,
			y:    2,
		},
		{
			name: "neg pos",
			x:    -1,
			y:    2,
		},
		{
			name: "neg neg",
			x:    -1,
			y:    -2,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			xBytes := int64ToBytes(c.x)
			yBytes := int64ToBytes(c.y)
			require.Equal(t, c.x > c.y, bytes.Compare(xBytes, yBytes) > 0)
			require.Equal(t, c.x == c.y, bytes.Equal(xBytes, yBytes))
		})

	}
}

func FuzzInt64Bytes_Order(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y int64) {
		xBytes := int64ToBytes(x)
		yBytes := int64ToBytes(y)
		require.Equal(t, x > y, bytes.Compare(xBytes, yBytes) > 0)
		require.Equal(t, x == y, bytes.Equal(xBytes, yBytes))
	})
}

func FuzzInt64Bytes_Float64(f *testing.F) {
	f.Fuzz(func(t *testing.T, x int64) {
		require.Equal(t, NewInt64(x).Int64(), x)
	})
}

func TestFloat64Bytes_Order(t *testing.T) {

	for _, c := range []struct {
		name string
		x    float64
		y    float64
	}{
		{
			name: "pos pos",
			x:    1,
			y:    2,
		},
		{
			name: "neg pos",
			x:    -1,
			y:    2,
		},
		{
			name: "neg neg",
			x:    -1,
			y:    -2,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			xBytes := float64ToBytes(c.x)
			yBytes := float64ToBytes(c.y)
			require.Equal(t, c.x > c.y, bytes.Compare(xBytes, yBytes) > 0)
			require.Equal(t, c.x == c.y, bytes.Equal(xBytes, yBytes))
		})

	}
}

func FuzzFloat64Bytes_Order(f *testing.F) {
	f.Fuzz(func(t *testing.T, x, y float64) {
		xBytes := float64ToBytes(x)
		yBytes := float64ToBytes(y)
		require.Equal(t, x > y, bytes.Compare(xBytes, yBytes) > 0)
		require.Equal(t, x == y, bytes.Equal(xBytes, yBytes))
	})
}

func FuzzFloat64Bytes_Float64(f *testing.F) {
	f.Fuzz(func(t *testing.T, x float64) {
		require.Equal(t, NewFloat64(x).Float64(), x)
	})
}

func FuzzSignFlip(f *testing.F) {
	f.Fuzz(func(t *testing.T, x float64) {
		require.Equal(t, math.Float64frombits(signFlip(math.Float64bits(x))), -x)
	})
}
