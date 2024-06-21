package api

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"

	"github.com/jamestunnell/marketanalysis/app/backend/background"
	"github.com/jamestunnell/marketanalysis/app/backend/stores"
	"github.com/jamestunnell/marketanalysis/graph"
)

type Graphs struct {
	*CRUDAPI[*graph.Config]
	DB *mongo.Database
	BG background.System
}

func NewGraphs(db *mongo.Database, bg background.System) *Graphs {
	return &Graphs{
		CRUDAPI: NewCRUDAPI(stores.NewGraphs(db)),
		DB:      db,
		BG:      bg,
	}
}

func (a *Graphs) Bind(r *mux.Router) {
	a.CRUDAPI.Bind(r)

	r.HandleFunc(a.PluralRoute()+"/run", a.Run).Methods(http.MethodPost)           //, http.MethodOptions)
	r.HandleFunc(a.PluralRoute()+"/optimize", a.Optimize).Methods(http.MethodPost) //, http.MethodOptions)
}
