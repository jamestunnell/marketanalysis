package graph_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	"github.com/jamestunnell/marketanalysis/blocks"
// )

// func TestJSONHappyPath(t *testing.T) {
// 	for _, typ := range blocks.Registry().Types() {
// 		testJSONHappyPath(t, typ)
// 	}
// }

// func testJSONHappyPath(t *testing.T, typ string) {
// 	t.Run(typ, func(t *testing.T) {
// 		newFunc, found := blocks.Registry().Get(typ)

// 		require.True(t, found)

// 		blk := newFunc()

// 		d, err := blocks.MarshalBlockJSON(blk)

// 		require.NoError(t, err)

// 		blk2, err := blocks.UnmarshalBlockJSON(d)

// 		require.NoError(t, err)

// 		assert.Equal(t, blk.GetType(), blk2.GetType())
// 	})
// }
