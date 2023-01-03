package term

import (
	"bytes"
	"fmt"
	"log"
)

type Term []byte

func (t Term) Raw() []byte {
	return t
}

func (t Term) Text() string {
	return string(t[1:])
}

func (t Term) Int64() int64 {
	return newInt64BytesFromBytes(t[1:]).int64()
}

func (t Term) Float64() float64 {
	return newFloat64BytesFromBytes(t[1:]).float64()
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
