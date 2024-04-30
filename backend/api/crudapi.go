package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type CRUDAPI[T any] struct {
	*Resource[T]

	Collection *mongo.Collection
}

func NewCRUDAPI[T any](r *Resource[T], db *mongo.Database) *CRUDAPI[T] {
	return &CRUDAPI[T]{
		Resource:   r,
		Collection: db.Collection(r.NamePlural),
	}
}

func (a *CRUDAPI[T]) Bind(r *mux.Router) {
	singleRoute := a.SingleRoute()

	r.HandleFunc(singleRoute, a.Put).Methods(http.MethodPut)
	r.HandleFunc(singleRoute, a.Get).Methods(http.MethodGet)
	r.HandleFunc(singleRoute, a.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/"+a.NamePlural, a.GetAll).Methods(http.MethodGet)
}

func (a *CRUDAPI[T]) Get(w http.ResponseWriter, r *http.Request) {
	Get(w, r, a.Resource, a.Collection)
}

func (a *CRUDAPI[T]) GetAll(w http.ResponseWriter, r *http.Request) {
	GetAll[T](w, r, a.Resource, a.Collection)
}

func (a *CRUDAPI[T]) Put(w http.ResponseWriter, r *http.Request) {
	Put(w, r, a.Resource, a.Collection)
}

func (a *CRUDAPI[T]) Delete(w http.ResponseWriter, r *http.Request) {
	Delete(w, r, a.Resource, a.Collection)
}

func (a *CRUDAPI[T]) FindOne(ctx context.Context, key string) (*T, *HTTPErr) {
	return FindOne(ctx, key, a.Resource, a.Collection)
}

func (a *CRUDAPI[T]) SingleRoute() string {
	return fmt.Sprintf("/%s/{%s}", a.Resource.NamePlural, a.Resource.KeyName)
}
