package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
	"go.mongodb.org/mongo-driver/mongo"
)

type API[T any] struct {
	KeyName    string
	Name       string
	NamePlural string
	Schema     *gojsonschema.Schema
	Collection *mongo.Collection
	Validate   func(*T) error
}

func (a *API[T]) Bind(r *mux.Router) {
	singleRoute := fmt.Sprintf("/%s/{%s}", a.NamePlural, a.KeyName)

	r.HandleFunc(singleRoute, a.Put).Methods(http.MethodPut)
	r.HandleFunc(singleRoute, a.Get).Methods(http.MethodGet)
	r.HandleFunc(singleRoute, a.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/"+a.NamePlural, a.GetAll).Methods(http.MethodGet)
}
