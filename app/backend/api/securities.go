package api

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app/stores"
	"github.com/jamestunnell/marketanalysis/models"
)

func NewSecurities(db *mongo.Database) *CRUDAPI[*models.Security] {
	return NewCRUDAPI(stores.NewSecurities(db))
}
