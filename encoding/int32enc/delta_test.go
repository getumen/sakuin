package int32enc_test

import (
	"testing"

	"github.com/getumen/sakuin/encoding/int32enc"
	"github.com/stretchr/testify/require"
)

func TestDelta(t *testing.T) {
	expected := []uint32{1, 3, 5, 7, 9}
	actual := make([]uint32, len(expected))
	copy(actual, expected)
	int32enc.EncodeDeltaInplace(actual)
	int32enc.DecodeDeltaInplace(actual)
	require.Equal(t, expected, actual)
}
