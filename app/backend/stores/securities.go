package stores

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/models"
)

func NewSecurities(db *mongo.Database) backend.Store[*models.Security] {
	info := &backend.ResourceInfo{
		KeyName:    models.SecurityKeyName,
		Name:       "security",
		NamePlural: "securities",
	}
	col := db.Collection(info.NamePlural)

	return backend.NewMongoStore[*models.Security](info, col)
}
