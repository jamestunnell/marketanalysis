package stores

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/graph"
)

func NewGraphs(db *mongo.Database) backend.Store[*graph.Configuration] {
	info := &backend.ResourceInfo{
		KeyName:    graph.ConfigKeyName,
		Name:       "graph",
		NamePlural: "graphs",
	}
	col := db.Collection(info.NamePlural)

	return backend.NewMongoStore[*graph.Configuration](info, col)
}
