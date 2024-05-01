package api

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"

	"github.com/jamestunnell/marketanalysis/app"
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
	rdef := &app.ResourceDef[graph.Configuration]{
		KeyName:    GraphKeyName,
		Name:       "graph",
		NamePlural: "graphs",
		Validate: func(cfg *graph.Configuration) []error {
			return cfg.Validate()
		},
		GetKey: func(c *graph.Configuration) string {
			return c.ID
		},
	}
	col := db.Collection(rdef.NamePlural)
	store := app.NewMongoStore[graph.Configuration](rdef, col)
	graphs := &Graphs{
		CRUDAPI:    NewCRUDAPI[graph.Configuration](store),
		securities: securities,
	}

	return graphs, nil
}

func (a *Graphs) Bind(r *mux.Router) {
	a.CRUDAPI.Bind(r)

	r.HandleFunc(a.SingularRoute()+"/run-day", a.RunDay).Methods(http.MethodPost)
}
