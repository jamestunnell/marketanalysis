package stores

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/models"
)

func NewBarSets(symbol string, db *mongo.Database) backend.Store[*models.BarSet] {
	info := &backend.ResourceInfo{
		KeyName:    "date",
		Name:       "barset",
		NamePlural: "barsets",
	}
	col := db.Collection(symbol)

	return backend.NewMongoStore[*models.BarSet](info, col)
}
