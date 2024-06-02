package api

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"

	"github.com/jamestunnell/marketanalysis/app/stores"
	"github.com/jamestunnell/marketanalysis/graph"
)

type Graphs struct {
	*CRUDAPI[*graph.Configuration]
}

func NewGraphs(
	db *mongo.Database,
) *Graphs {
	return &Graphs{
		CRUDAPI: NewCRUDAPI(stores.NewGraphs(db)),
	}
}

func (a *Graphs) Bind(r *mux.Router) {
	a.CRUDAPI.Bind(r)

	r.HandleFunc(a.SingularRoute()+"/backtest", a.BacktestGraph).Methods(http.MethodPost) //, http.MethodOptions)
	r.HandleFunc(a.SingularRoute()+"/eval", a.EvalGraph).Methods(http.MethodPost)         //, http.MethodOptions)
	r.HandleFunc(a.SingularRoute()+"/run", a.RunGraph).Methods(http.MethodPost)           //, http.MethodOptions)
}
