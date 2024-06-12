package stores

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app/backend"
	bemodels "github.com/jamestunnell/marketanalysis/app/backend/models"
)

func NewSettings(db *mongo.Database) backend.Store[*bemodels.Setting] {
	info := &backend.ResourceInfo{
		KeyName:    bemodels.SettingKeyName,
		Name:       "setting",
		NamePlural: "settings",
	}
	col := db.Collection(info.NamePlural)

	return backend.NewMongoStore[*bemodels.Setting](info, col)
}
