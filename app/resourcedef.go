package api

import (
	"github.com/xeipuuv/gojsonschema"
)

type Resource[T any] struct {
	KeyName    string
	Name       string
	NamePlural string
	Schema     *gojsonschema.Schema
	Validate   func(*T) error
}
