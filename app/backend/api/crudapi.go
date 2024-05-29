package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app"
)

type CRUDAPI[T app.Resource] struct {
	Store app.Store[T]
}

func NewCRUDAPI[T app.Resource](s app.Store[T]) *CRUDAPI[T] {
	return &CRUDAPI[T]{
		Store: s,
	}
}

func (a *CRUDAPI[T]) Bind(r *mux.Router) {
	r.HandleFunc(a.PluralRoute(), a.handlePlural).
		Methods(http.MethodGet, http.MethodPost) //, http.MethodOptions)

	r.HandleFunc(a.SingularRoute(), a.handleSingle).
		Methods(http.MethodGet, http.MethodPut, http.MethodDelete) //, http.MethodOptions)
}

func (a *CRUDAPI[T]) Get(w http.ResponseWriter, r *http.Request) {
	Get(w, r, a.Store)
}

func (a *CRUDAPI[T]) handlePlural(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.GetAll(w, r)
	case http.MethodPost:
		a.Create(w, r)
	default:
		log.Error().Msgf("unexpected HTTP method %s", r.Method)
	}
}

func (a *CRUDAPI[T]) handleSingle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.Get(w, r)
	case http.MethodPut:
		a.Update(w, r)
	case http.MethodDelete:
		a.Delete(w, r)
	default:
		log.Error().Msgf("unexpected HTTP method %s", r.Method)
	}
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
	return fmt.Sprintf("/%s", a.Store.GetInfo().NamePlural)
}

func (a *CRUDAPI[T]) SingularRoute() string {
	return fmt.Sprintf("/%s/{%s}", a.Store.GetInfo().NamePlural, a.Store.GetInfo().KeyName)
}
