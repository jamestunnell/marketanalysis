package api

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"

	"github.com/jamestunnell/marketanalysis/app/stores"
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/models"
)

type Graphs struct {
	*CRUDAPI[*graph.Configuration]

	securities *CRUDAPI[*models.Security]
}

func NewGraphs(
	db *mongo.Database,
	securities *CRUDAPI[*models.Security],
) *Graphs {
	return &Graphs{
		CRUDAPI:    NewCRUDAPI(stores.NewGraphs(db)),
		securities: securities,
	}
}

func (a *Graphs) Bind(r *mux.Router) {
	a.CRUDAPI.Bind(r)

	r.HandleFunc(a.SingularRoute()+"/eval", a.EvalGraph).Methods(http.MethodPost) //, http.MethodOptions)
	r.HandleFunc(a.SingularRoute()+"/run", a.RunGraph).Methods(http.MethodPost)   //, http.MethodOptions)
}
