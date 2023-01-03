package int32enc

import (
	"encoding/binary"
	"fmt"

	"github.com/bmkessler/streamvbyte"
)

const (
	raw uint64 = iota
	streamVByte
)

func Encode(values []uint32) ([]byte, error) {
	result := make([]byte, 0)
	result = binary.AppendUvarint(result, streamVByte)
	blob, err := encodeSteamVbyte(values)
	if err != nil {
		return nil, err
	}
	result = binary.AppendUvarint(result, uint64(len(blob)))
	result = append(result, blob...)
	return result, nil
}

func Decode(blob []byte) (values []uint32, err error) {
	var cur int
	enc, diff := binary.Uvarint(blob[cur:])
	cur += diff
	blobSize, diff := binary.Uvarint(blob[cur:])
	cur += diff
	switch enc {
	case streamVByte:
		values, err = decodeSteamVbyte(blob[cur : cur+int(blobSize)])
		return values, err
	default:
		return nil, fmt.Errorf("unknown encoding: %d", enc)
	}
}

func encodeSteamVbyte(values []uint32) ([]byte, error) {
	result := make([]byte, 0)
	result = binary.AppendUvarint(result, uint64(len(values)))

	encoded := make([]byte, streamvbyte.MaxSize32(len(values)))
	n := streamvbyte.EncodeUint32(encoded, values)
	result = append(result, encoded[:n]...)
	return result, nil
}

func decodeSteamVbyte(blob []byte) ([]uint32, error) {
	var cur int
	valueNum, diff := binary.Uvarint(blob[cur:])
	cur += diff

	result := make([]uint32, valueNum)
	streamvbyte.DecodeUint32(result, blob[cur:])
	return result, nil
}
