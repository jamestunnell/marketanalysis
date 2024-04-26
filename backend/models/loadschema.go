package models

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func LoadSchema(str string) (*gojsonschema.Schema, error) {
	l := gojsonschema.NewStringLoader(str)

	schema, err := gojsonschema.NewSchema(l)
	if err != nil {
		err = fmt.Errorf("failed to make JSON schema: %w", err)

		return nil, err
	}

	return schema, nil
}
