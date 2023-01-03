package int64enc_test

import (
	"testing"

	"github.com/getumen/sakuin/encoding/int64enc"
	"github.com/stretchr/testify/require"
)

func TestDelta(t *testing.T) {
	expected := []uint64{1, 3, 5, 7, 9}
	actual := make([]uint64, len(expected))
	copy(actual, expected)
	int64enc.EncodeDeltaInplace(actual)
	int64enc.DecodeDeltaInplace(actual)
	require.Equal(t, expected, actual)
}
