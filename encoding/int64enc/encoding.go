package int64enc

import (
	"encoding/binary"
	"fmt"

	"github.com/jwilder/encoding/simple8b"
)

const (
	raw uint64 = iota
	simple8B
)

func Encode(values []uint64) ([]byte, error) {
	result := make([]byte, 0)
	result = binary.AppendUvarint(result, simple8B)
	blob, err := encodeSimple8B(values)
	if err != nil {
		return nil, err
	}
	result = binary.AppendUvarint(result, uint64(len(blob)))
	result = append(result, blob...)
	return result, nil
}

func Decode(blob []byte) (values []uint64, err error) {
	var cur int
	enc, diff := binary.Uvarint(blob[cur:])
	cur += diff
	blobSize, diff := binary.Uvarint(blob[cur:])
	cur += diff
	switch enc {
	case simple8B:
		values, err = decodeSimple8B(blob[cur : cur+int(blobSize)])
		return values, err
	default:
		return nil, fmt.Errorf("unknown encoding: %d", enc)
	}
}

func encodeSimple8B(values []uint64) ([]byte, error) {
	encoder := simple8b.NewEncoder()

	encoder.SetValues(values)

	return encoder.Bytes()
}

func decodeSimple8B(blob []byte) ([]uint64, error) {
	result := make([]uint64, 0)

	decoder := simple8b.NewDecoder(blob)
	for decoder.Next() {
		result = append(result, decoder.Read())
	}
	return result, nil
}
