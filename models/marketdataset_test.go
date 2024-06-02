package models_test

import (
	"testing"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/stretchr/testify/assert"
)

const testSymbol = "TST"

func TestMarketDataset(t *testing.T) {
	ds := models.NewMarketDataset(testSymbol)

	assert.Empty(t, ds.Bars)
}
