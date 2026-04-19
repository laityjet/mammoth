package stream

import (
	"testing"

	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/require"
)

func TestMerger(t *testing.T) {
	its := []Peekable[string]{
		NewSlice([]string{"c", "d"}),
		NewSlice([]string{"a"}),
		NewSlice([]string{"b", "c"}),
		NewSlice([]string{}),
	}

	expected := []Merged[string]{
		{Indexes: []int{1}, Values: []string{"a"}},
		{Indexes: []int{2}, Values: []string{"b"}},
		{Indexes: []int{0, 2}, Values: []string{"c", "c"}},
		{Indexes: []int{0}, Values: []string{"d"}},
	}
	m := NewMerger(its, func(a, b string) bool {
		return a < b
	})

	actual, err := Collect[Merged[string]](pctx.TestContext(t), m)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
