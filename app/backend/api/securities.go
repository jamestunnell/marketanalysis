package api

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"

	"github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/app/stores"
)

type Securities struct {
	*CRUDAPI[*models.Security]
	DB *mongo.Database
}

func NewSecurities(db *mongo.Database) *Securities {
	return &Securities{
		CRUDAPI: NewCRUDAPI(stores.NewSecurities(db)),
		DB:      db,
	}
}

func (a *Securities) Bind(r *mux.Router) {
	a.CRUDAPI.Bind(r)

	// r.HandleFunc(a.SingularRoute()+"/collect", a.BacktestGraph).Methods(http.MethodPost) //, http.MethodOptions)
	// r.HandleFunc(a.SingularRoute()+"/eval", a.EvalGraph).Methods(http.MethodPost)        //, http.MethodOptions)
	// r.HandleFunc(a.SingularRoute()+"/run", a.RunGraph).Methods(http.MethodPost)          //, http.MethodOptions)
}
