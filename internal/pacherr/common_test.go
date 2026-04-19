package pacherr

import (
	"testing"

	"github.com/laityjet/mammoth/v0/internal/require"
)

func TestIsNotExist(t *testing.T) {
	err := NewNotExist("collection", "id")
	require.True(t, IsNotExist(err))
}

func TestIsExist(t *testing.T) {
	err := NewExists("collection", "id")
	require.True(t, IsExists(err))
}
