package dispatch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookingUpRecentValues(t *testing.T) {
	require := require.New(t)
	r := newRecent(2)
	require.False(r.lookup("msg-id"))
	r.insert("msg-id")
	require.True(r.lookup("msg-id"))
}

func TestLookingUpValuesThatHaveWerePushedOut(t *testing.T) {
	require := require.New(t)
	r := newRecent(2)
	r.insert("msg-id-1")
	r.insert("msg-id-2")
	require.True(r.lookup("msg-id-1"))
	r.insert("msg-id-3")
	require.False(r.lookup("msg-id-1"))
}
