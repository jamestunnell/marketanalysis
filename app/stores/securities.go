package stores

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/app/backend/models"
)

func NewSecurities(db *mongo.Database) app.Store[*models.Security] {
	info := &app.ResourceInfo{
		KeyName:    models.SecurityKeyName,
		Name:       "security",
		NamePlural: "securities",
	}
	col := db.Collection(info.NamePlural)

	return app.NewMongoStore[*models.Security](info, col)
}
