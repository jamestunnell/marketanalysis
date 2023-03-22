package datamanager_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jamestunnell/marketanalysis/datamanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectionInfo_StoreLoad(t *testing.T) {
	info := datamanager.NewCollectionInfo(
		"QQQ", datamanager.Resolution1Min)

	dir, err := os.MkdirTemp("", "dmtests*")

	require.NoError(t, err)

	fpath := filepath.Join(dir, "collectioninfo.json")

	err = info.Store(fpath)

	require.NoError(t, err)

	var info2 datamanager.CollectionInfo

	err = (&info2).Load(fpath)

	require.NoError(t, err)

	assert.Equal(t, info2.Symbol, info.Symbol)
	assert.Equal(t, info2.Resolution, info.Resolution)
}
