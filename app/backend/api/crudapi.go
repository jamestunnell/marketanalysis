package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/rs/zerolog/log"
)

type CRUDAPI[T backend.Resource] struct {
	Store   backend.Store[T]
	Options *CRUDOptions[T]
}

type CRUDOptions[T backend.Resource] struct {
	Hooks *CRUDHooks[T]
}

type CRUDHooks[T backend.Resource] struct {
	PostCreate func(res T)
	PostUpdate func(res T)
	PostDelete func(key string)
}

type CRUDOptionsMod[T backend.Resource] func(*CRUDOptions[T])

func WithPostCreateHook[T backend.Resource](hook func(res T)) CRUDOptionsMod[T] {
	return func(opts *CRUDOptions[T]) {
		opts.Hooks.PostCreate = hook
	}
}

func WithPostUpdateHook[T backend.Resource](hook func(res T)) CRUDOptionsMod[T] {
	return func(opts *CRUDOptions[T]) {
		opts.Hooks.PostUpdate = hook
	}
}

func WithPostDeleteHook[T backend.Resource](hook func(key string)) CRUDOptionsMod[T] {
	return func(opts *CRUDOptions[T]) {
		opts.Hooks.PostDelete = hook
	}
}

func NewCRUDAPI[T backend.Resource](s backend.Store[T], mods ...CRUDOptionsMod[T]) *CRUDAPI[T] {
	opts := &CRUDOptions[T]{
		Hooks: &CRUDHooks[T]{},
	}

	for _, mod := range mods {
		mod(opts)
	}

	return &CRUDAPI[T]{
		Store:   s,
		Options: opts,
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
		a.Upsert(w, r)
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
	Create(w, r, a.Store, a.Options.Hooks.PostCreate)
}

func (a *CRUDAPI[T]) Upsert(w http.ResponseWriter, r *http.Request) {
	Upsert(w, r, a.Store, a.Options.Hooks.PostUpdate)
}

func (a *CRUDAPI[T]) Delete(w http.ResponseWriter, r *http.Request) {
	Delete(w, r, a.Store, a.Options.Hooks.PostDelete)
}

func (a *CRUDAPI[T]) PluralRoute() string {
	return fmt.Sprintf("/%s", a.Store.GetInfo().NamePlural)
}

func (a *CRUDAPI[T]) SingularRoute() string {
	return fmt.Sprintf("/%s/{%s}", a.Store.GetInfo().NamePlural, a.Store.GetInfo().KeyName)
}
