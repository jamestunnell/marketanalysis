package api

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/models"
)

func NewSecurities(db *mongo.Database) (*CRUDAPI[models.Security], error) {
	rdef := &app.ResourceDef[models.Security]{
		KeyName:    models.SecurityKeyName,
		Name:       models.SecurityName,
		NamePlural: models.SecurityNamePlural,
		Validate: func(s *models.Security) []error {
			return s.Validate()
		},
		GetKey: func(s *models.Security) string {
			return s.Symbol
		},
	}
	col := db.Collection(rdef.NamePlural)
	store := app.NewMongoStore[models.Security](rdef, col)

	a := NewCRUDAPI[models.Security](store)

	return a, nil
}
