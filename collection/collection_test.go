package collection_test

import (
	"os"
	"testing"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollection(t *testing.T) {
	store, cleanup := makeTestStore(t)
	defer cleanup()

	bars := makeTestBars(t)
	info := &collection.Info{
		Symbol:     "QQQ",
		Resolution: collection.Resolution1Min,
	}

	c, err := collection.New(info, bars)

	assert.NoError(t, err)

	err = c.Store(store)

	assert.NoError(t, err)
}

func makeTestStore(t *testing.T) (s collection.Store, cleanup func()) {
	root, err := os.MkdirTemp("", "collectiontests")

	require.NoError(t, err)

	cleanup = func() {
		os.RemoveAll(root)
	}

	store, err := collection.NewDirStore(root)

	if !assert.NoError(t, err) {
		cleanup()

		t.Fatalf("failed to make dir store: %v", err)
	}

	return store, cleanup
}

func makeTestBars(t *testing.T) []*models.Bar {
	const testBarsJSON = `
{"t":"2023-03-16T13:30:00Z","o":386.82,"h":387.32,"l":386.72,"c":386.72,"v":725510,"n":5656,"vw":386.97766}
{"t":"2023-03-16T13:31:00Z","o":386.72,"h":386.89,"l":386.5,"c":386.56,"v":481409,"n":4916,"vw":386.70493}
{"t":"2023-03-16T13:32:00Z","o":386.55,"h":386.76,"l":386.29,"c":386.4611,"v":562639,"n":5215,"vw":386.51794}
{"t":"2023-03-16T13:33:00Z","o":386.45,"h":386.82,"l":386.39,"c":386.62,"v":398813,"n":3774,"vw":386.606}
{"t":"2023-03-16T13:34:00Z","o":386.64,"h":387.18,"l":386.48,"c":387.099,"v":444607,"n":4265,"vw":386.76712}
`

	f, err := os.CreateTemp("", "collectiontests*")

	require.NoError(t, err)

	fname := f.Name()
	defer os.Remove(fname)

	_, err = f.WriteString(testBarsJSON)

	require.NoError(t, err)

	require.NoError(t, f.Close())

	bars, err := models.LoadBars(fname)

	require.NoError(t, err)

	return bars
}
