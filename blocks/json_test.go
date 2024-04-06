package blocks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jamestunnell/marketanalysis/blocks"
)

func TestJSONHappyPath(t *testing.T) {
	testJSONHappyPath(t, blocks.TypeSMA)
	testJSONHappyPath(t, blocks.TypeEMA)
}

func testJSONHappyPath(t *testing.T, typ string) {
	newFunc, found := blocks.Registry().Get(typ)

	require.True(t, found)

	blk := newFunc()

	d, err := blocks.MarshalJSON(blk)

	require.NoError(t, err)

	blk2, err := blocks.UnmarshalJSON(d)

	require.NoError(t, err)

	assert.Equal(t, blk.GetType(), blk2.GetType())
}
