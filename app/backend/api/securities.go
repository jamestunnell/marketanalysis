package api

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/models"
)

func NewSecurities(db *mongo.Database) (*CRUDAPI[*models.Security], error) {
	info := &app.ResourceInfo{
		KeyName:    models.SecurityKeyName,
		Name:       models.SecurityName,
		NamePlural: models.SecurityNamePlural,
	}
	col := db.Collection(info.NamePlural)
	store := app.NewMongoStore[*models.Security](info, col)

	a := NewCRUDAPI(store)

	return a, nil
}
