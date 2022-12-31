package term

import (
	"bytes"
	"fmt"
	"log"
)

type Term []byte

func (t Term) String() string {
	switch t[0] {
	case byte(Nil):
		return "nil"
	case byte(Text):
		return string(t[1:])
	case byte(Int64):
		return fmt.Sprintf("%d", newInt64BytesFromBytes(t[1:]).int64())
	case byte(Float64):
		return fmt.Sprintf("%f", newFloat64BytesFromBytes(t[1:]).float64())
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
	b = append(b, byte(Text))
	b = append(b, newInt64Bytes(value).bytes()...)
	return b
}

func NewFloat64(value float64) Term {
	b := make([]byte, 0)
	b = append(b, byte(Text))
	b = append(b, newFloat64Bytes(value).bytes()...)
	return b
}
