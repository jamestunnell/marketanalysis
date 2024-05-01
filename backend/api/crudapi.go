package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jamestunnell/marketanalysis/app"
)

type CRUDAPI[T any] struct {
	Store app.Store[T]
}

func NewCRUDAPI[T any](s app.Store[T]) *CRUDAPI[T] {
	return &CRUDAPI[T]{
		Store: s,
	}
}

func (a *CRUDAPI[T]) Bind(r *mux.Router) {
	r.HandleFunc(a.PluralRoute(), a.GetAll).Methods(http.MethodGet)
	r.HandleFunc(a.PluralRoute(), a.Create).Methods(http.MethodPost)

	r.HandleFunc(a.SingularRoute(), a.Get).Methods(http.MethodGet)
	r.HandleFunc(a.SingularRoute(), a.Update).Methods(http.MethodPut)
	r.HandleFunc(a.SingularRoute(), a.Delete).Methods(http.MethodDelete)
}

func (a *CRUDAPI[T]) Get(w http.ResponseWriter, r *http.Request) {
	Get(w, r, a.Store)
}

func (a *CRUDAPI[T]) GetAll(w http.ResponseWriter, r *http.Request) {
	GetAll[T](w, r, a.Store)
}

func (a *CRUDAPI[T]) Create(w http.ResponseWriter, r *http.Request) {
	Create(w, r, a.Store)
}

func (a *CRUDAPI[T]) Update(w http.ResponseWriter, r *http.Request) {
	Update(w, r, a.Store)
}

func (a *CRUDAPI[T]) Delete(w http.ResponseWriter, r *http.Request) {
	Delete(w, r, a.Store)
}

func (a *CRUDAPI[T]) PluralRoute() string {
	return fmt.Sprintf("/%s", a.Store.RDef().NamePlural)
}

func (a *CRUDAPI[T]) SingularRoute() string {
	return fmt.Sprintf("/%s/{%s}", a.Store.RDef().NamePlural, a.Store.RDef().KeyName)
}
