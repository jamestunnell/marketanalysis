package collection_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirStore(t *testing.T) {
	root, err := os.MkdirTemp("", "collectiontests")

	require.NoError(t, err)

	defer os.RemoveAll(root)

	store, err := collection.NewDirStore(root)

	require.NoError(t, err)

	// nothing in the store yet
	assert.NoError(t, err)
	assert.Empty(t, store.ItemNames())

	// load non-existent item fails
	d, err := store.LoadItem("bogus")

	assert.Empty(t, d)
	assert.Error(t, err)

	info := collection.NewInfo(
		"QQQ", collection.Resolution1Min)

	d, err = json.Marshal(info)

	require.NoError(t, err)

	// store a new item succeeds
	err = store.StoreItem("info.json", d)

	require.NoError(t, err)

	// one item in the store
	names := store.ItemNames()

	assert.NoError(t, err)
	assert.Len(t, names, 1)
	assert.Contains(t, names, "info.json")

	// load existent item succeeds
	d2, err := store.LoadItem("info.json")

	require.NoError(t, err)

	assert.Equal(t, d, d2)
}
