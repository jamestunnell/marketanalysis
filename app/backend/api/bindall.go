package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jamestunnell/marketanalysis/app/backend/stores"
	"go.mongodb.org/mongo-driver/mongo"
)

func BindAll(r *mux.Router, db *mongo.Database) {
	graphs := NewGraphs(db)
	settings := NewCRUDAPI(stores.NewSettings(db))

	r.Handle("/blocks/{type}", NewGetBlockInfo()).Methods(http.MethodGet) //, http.MethodOptions)
	r.Handle("/blocks", NewGetAllBlockInfo()).Methods(http.MethodGet)     //, http.MethodOptions)
	r.Handle("/status", NewStatus()).Methods(http.MethodPost)             //, http.MethodOptions)

	r.Handle("/bars/{symbol}/{date}", NewGetBars(db))

	settings.Bind(r)
	graphs.Bind(r)
}
