package term

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
)

type Term []byte

func (t Term) Raw() []byte {
	return t
}

func (t Term) Text() string {
	return string(t[1:])
}

func (t Term) Int64() int64 {
	return bytesToInt64(t[1:])
}

func (t Term) Float64() float64 {
	return bytesToFloat64(t[1:])
}

func (t Term) String() string {
	switch t[0] {
	case byte(Nil):
		return "nil"
	case byte(Text):
		return t.Text()
	case byte(Int64):
		return fmt.Sprintf("%d", t.Int64())
	case byte(Float64):
		return fmt.Sprintf("%f", t.Float64())
	default:
		log.Panicf("unsupported term type: %d", t[0])
		return ""
	}
}

func (t Term) TermType() TermType {
	switch t[0] {
	case byte(Nil):
		return Nil
	case byte(Text):
		return Text
	case byte(Int64):
		return Int64
	case byte(Float64):
		return Float64
	default:
		log.Panicf("unsupported term type: %d", t[0])
		return Nil
	}
}

func (t Term) Type() string {
	switch t[0] {
	case byte(Nil):
		return "nil"
	case byte(Text):
		return "text"
	case byte(Int64):
		return "int64"
	case byte(Float64):
		return "float64"
	default:
		log.Panicf("unsupported term type: %d", t[0])
		return ""
	}
}

type TermType byte

func Comparator(x, y interface{}) int {
	return bytes.Compare(x.(Term), y.(Term))
}

// don't modify the order
const (
	Nil TermType = iota
	Text
	Int64
	Float64
)

func NewText(value string) Term {
	b := make([]byte, 0)
	b = append(b, byte(Text))
	b = append(b, []byte(value)...)
	return b
}

func NewInt64(value int64) Term {
	b := make([]byte, 0)
	b = append(b, byte(Int64))
	b = append(b, int64ToBytes(value)...)
	return b
}

func int64ToBytes(value int64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, uint64(value)^(1<<63))
	return result
}

func bytesToInt64(bytes []byte) int64 {
	return int64(binary.BigEndian.Uint64(bytes) ^ (1 << 63))
}

func NewFloat64(value float64) Term {
	b := make([]byte, 0)
	b = append(b, byte(Float64))
	b = append(b, float64ToBytes(value)...)
	return b
}

func float64ToBytes(value float64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, signFlip(math.Float64bits(value)))
	return result
}

func bytesToFloat64(bytes []byte) float64 {
	return math.Float64frombits(signFlip(binary.BigEndian.Uint64(bytes)))
}

// https://lemire.me/blog/2020/12/14/converting-floating-point-numbers-to-integers-while-preserving-order/?utm_source=pocket_saves
func signFlip(x uint64) uint64 {
	mask := uint64(int64(x) >> 63)
	mask |= 0x8000000000000000
	return (x ^ mask)
}
