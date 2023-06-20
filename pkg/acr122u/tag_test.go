package acr122u

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
	// 03 18 d101145504737369732e636f746e6574776f726b2e636f6dfe
	// d1
	// 1011 0001
	b, err := hex.DecodeString("0318d1011355046469643a3a31323334353637383930fe")
	require.NoError(t, err)
	tag := NewTag(nil, "")
	err = tag.Unmarshal(b)
	require.NoError(t, err)
	t.Logf("%s\n", tag.Message.Payload)
	t.Logf("%#v\n", tag)
	t.Logf("%#v\n", tag.Message)

	bb, err := tag.Marshal()
	require.NoError(t, err)
	require.Equal(t, b, bb)
}
