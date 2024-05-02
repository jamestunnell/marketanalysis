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
	*CRUDAPI[*graph.Configuration]

	securities *CRUDAPI[*models.Security]
}

func NewGraphs(
	db *mongo.Database,
	securities *CRUDAPI[*models.Security],
) (*Graphs, error) {
	info := &app.ResourceInfo{
		KeyName:    graph.ConfigKeyName,
		Name:       "graph",
		NamePlural: "graphs",
	}
	col := db.Collection(info.NamePlural)
	store := app.NewMongoStore[*graph.Configuration](info, col)
	graphs := &Graphs{
		CRUDAPI:    NewCRUDAPI(store),
		securities: securities,
	}

	return graphs, nil
}

func (a *Graphs) Bind(r *mux.Router) {
	a.CRUDAPI.Bind(r)

	r.HandleFunc(a.SingularRoute()+"/run-day", a.RunDay).Methods(http.MethodPost)
}
