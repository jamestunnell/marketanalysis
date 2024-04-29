package api

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/backend/models"
)

func MakeSecuritiesAPI(db *mongo.Database) (*API[models.Security], error) {
	schema, err := models.LoadSchema(models.SecuritySchemaStr)
	if err != nil {
		return nil, fmt.Errorf("failed to load security schema: %w", err)
	}

	a := &API[models.Security]{
		KeyName:    models.SecurityKeyName,
		Name:       models.SecurityName,
		NamePlural: models.SecurityNamePlural,
		Collection: db.Collection(models.SecurityNamePlural),
		Schema:     schema,
		Validate: func(s *models.Security) error {
			if s.Open.MinuteOfDay() >= s.Close.MinuteOfDay() {
				return fmt.Errorf("open '%s' is not before close '%s'", s.Open, s.Close)
			}

			return nil
		},
	}

	return a, nil
}
