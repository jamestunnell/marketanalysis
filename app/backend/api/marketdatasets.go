package api

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"

	"github.com/jamestunnell/marketanalysis/app/stores"
	"github.com/jamestunnell/marketanalysis/models"
)

type MarketDatasets struct {
	*CRUDAPI[*models.MarketDataset]
}

func NewMarketDatasets(
	db *mongo.Database,
) *MarketDatasets {
	return &MarketDatasets{
		CRUDAPI: NewCRUDAPI(stores.NewMarketDatasets(db)),
	}
}

func (a *MarketDatasets) Bind(r *mux.Router) {
	a.CRUDAPI.Bind(r)
}
