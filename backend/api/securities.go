package api

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/models"
)

func NewSecurities(db *mongo.Database) (*CRUDAPI[models.Security], error) {
	schema, err := LoadSchema(models.SecuritySchemaStr)
	if err != nil {
		return nil, fmt.Errorf("failed to load security schema: %w", err)
	}

	res := &Resource[models.Security]{
		KeyName:    models.SecurityKeyName,
		Name:       models.SecurityName,
		NamePlural: models.SecurityNamePlural,
		Schema:     schema,
		Validate: func(s *models.Security) error {
			if s.Open.MinuteOfDay() >= s.Close.MinuteOfDay() {
				return fmt.Errorf("open '%s' is not before close '%s'", s.Open, s.Close)
			}

			if _, err := time.LoadLocation(s.TimeZone); err != nil {
				return fmt.Errorf("time zone '%s' is invalid: %w", err)
			}

			return nil
		},
	}

	a := NewCRUDAPI[models.Security](res, db)

	return a, nil
}
