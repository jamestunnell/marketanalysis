package api

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/jamestunnell/marketanalysis/app/backend/background"
	"github.com/jamestunnell/marketanalysis/app/backend/stores"
)

func BindAll(r *mux.Router, db *mongo.Database, bg background.System) {
	graphs := NewGraphs(db, bg)
	settings := NewCRUDAPI(stores.NewSettings(db))

	r.Handle("/blocks/{type}", NewGetBlockInfo()).Methods(http.MethodGet) //, http.MethodOptions)
	r.Handle("/blocks", NewGetAllBlockInfo()).Methods(http.MethodGet)     //, http.MethodOptions)
	r.Handle("/status", NewStatus()).Methods(http.MethodPost)             //, http.MethodOptions)

	r.Handle("/bars/{symbol}/{date}", NewGetBars(db))

	r.Handle("/jobs/{id}", NewJobStatus(bg))
	r.Handle("/jobs/updates", NewJobUpdates(bg))

	settings.Bind(r)
	graphs.Bind(r)
}
