package backend

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/backend/api"
)

func BindAPI(r *mux.Router, db *mongo.Database) {
	r.Handle("/status", api.NewStatus())

	db.Collection("assets").DeleteMany(context.Background(), bson.D{})

	coll := db.Collection("securities")

	r.Handle("/securities/{symbol}", api.NewPutSecurity(coll)).
		Methods(http.MethodPut)
	r.Handle("/securities/{symbol}", api.NewGetSecurity(coll)).
		Methods(http.MethodGet)
	r.Handle("/securities/{symbol}", api.NewDelSecurity(coll)).
		Methods(http.MethodDelete)
	r.Handle("/securities", api.NewGetSecurities(coll)).
		Methods(http.MethodGet)
}
