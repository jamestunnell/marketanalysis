package collection_test

import (
	"encoding/json"
	"testing"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfo_StoreLoad(t *testing.T) {
	info := &models.CollectionInfo{
		Symbol:     "QQQ",
		Resolution: models.Resolution1Min,
	}

	d, err := json.Marshal(info)

	require.NoError(t, err)

	var info2 models.CollectionInfo

	err = json.Unmarshal(d, &info2)

	assert.NoError(t, err)

	assert.Equal(t, info2.Symbol, info.Symbol)
	assert.Equal(t, info2.Resolution, info.Resolution)
}
