package collection_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfo_StoreLoad(t *testing.T) {
	info := collection.NewInfo(
		"QQQ", collection.Resolution1Min)

	dir, err := os.MkdirTemp("", "dmtests*")

	require.NoError(t, err)

	fpath := filepath.Join(dir, "Info.json")

	err = info.Store(fpath)

	require.NoError(t, err)

	var info2 collection.Info

	err = (&info2).Load(fpath)

	require.NoError(t, err)

	assert.Equal(t, info2.Symbol, info.Symbol)
	assert.Equal(t, info2.Resolution, info.Resolution)
}
