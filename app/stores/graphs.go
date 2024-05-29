package stores

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/graph"
)

func NewGraphs(db *mongo.Database) app.Store[*graph.Configuration] {
	info := &app.ResourceInfo{
		KeyName:    graph.ConfigKeyName,
		Name:       "graph",
		NamePlural: "graphs",
	}
	col := db.Collection(info.NamePlural)

	return app.NewMongoStore[*graph.Configuration](info, col)
}
