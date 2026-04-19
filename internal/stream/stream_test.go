package stream

import (
	"testing"

	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/require"
)

func TestIsEOS(t *testing.T) {
	require.True(t, IsEOS(EOS()))
	require.False(t, errors.Is(EOS(), EOS()))
}
