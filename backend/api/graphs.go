package api

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/models"
)

type Graphs struct {
	*CRUDAPI[graph.Configuration]

	securities *CRUDAPI[models.Security]
}

const GraphKeyName = "id"

func NewGraphs(
	db *mongo.Database,
	securities *CRUDAPI[models.Security],
) (*Graphs, error) {
	res := &Resource[graph.Configuration]{
		KeyName:    GraphKeyName,
		Name:       "graph",
		NamePlural: "graphs",
		Schema:     graph.GetConfigSchema(),
		Validate: func(cfg *graph.Configuration) error {
			return cfg.Validate()
		},
	}

	graphs := &Graphs{
		CRUDAPI:    NewCRUDAPI[graph.Configuration](res, db),
		securities: securities,
	}

	return graphs, nil
}

func (a *Graphs) Bind(r *mux.Router) {
	a.CRUDAPI.Bind(r)

	r.HandleFunc(a.SingleRoute(), a.Run).Methods(http.MethodPost)
}
