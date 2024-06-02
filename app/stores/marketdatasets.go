package stores

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/models"
)

const MarketDatasetKeyName = "symbol"

func NewMarketDatasets(db *mongo.Database) app.Store[*models.MarketDataset] {
	info := &app.ResourceInfo{
		KeyName:    MarketDatasetKeyName,
		Name:       "marketdataset",
		NamePlural: "marketdatasets",
	}
	col := db.Collection(info.NamePlural)

	return app.NewMongoStore[*models.MarketDataset](info, col)
}
