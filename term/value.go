package term

import (
	"encoding/binary"
	"math"
)

type int64Bytes uint64

func newInt64Bytes(value int64) int64Bytes {
	return int64Bytes(uint64(value) ^ (1 << 63))
}

func newInt64BytesFromBytes(b []byte) int64Bytes {
	return int64Bytes(binary.BigEndian.Uint64(b))
}

func (i int64Bytes) int64() int64 {
	return int64(uint64(i) ^ (1 << 63))
}

func (i int64Bytes) bytes() []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, uint64(i))
	return result
}

// https://lemire.me/blog/2020/12/14/converting-floating-point-numbers-to-integers-while-preserving-order/?utm_source=pocket_saves

type float64Bytes uint64

func newFloat64Bytes(value float64) float64Bytes {
	return float64Bytes(signFlip(math.Float64bits(value)))
}

func newFloat64BytesFromBytes(b []byte) float64Bytes {
	return float64Bytes(binary.BigEndian.Uint64(b))
}

func (f float64Bytes) float64() float64 {
	return math.Float64frombits(signFlip(uint64(f)))
}

func (f float64Bytes) bytes() []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, uint64(f))
	return result
}

func signFlip(x uint64) uint64 {
	mask := uint64(int64(x) >> 63)
	mask |= 0x8000000000000000
	return (x ^ mask)
}
