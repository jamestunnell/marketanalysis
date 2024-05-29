package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func BindAll(r *mux.Router, db *mongo.Database) {
	securities := NewSecurities(db)
	graphs := NewGraphs(db, securities)

	r.Handle("/blocks/{type}", NewGetBlockInfo()).Methods(http.MethodGet) //, http.MethodOptions)
	r.Handle("/blocks", NewGetAllBlockInfo()).Methods(http.MethodGet)     //, http.MethodOptions)
	r.Handle("/status", NewStatus()).Methods(http.MethodPost)             //, http.MethodOptions)

	graphs.Bind(r)
	securities.Bind(r)
}
