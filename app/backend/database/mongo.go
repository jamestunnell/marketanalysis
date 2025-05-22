package database

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/jamestunnell/marketanalysis/app/backend/appvars"
)

const MongoDBName = "marketanalysis"

func ConnectMongo(vars *appvars.Vars) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// var cred options.Credential

	// cred.AuthSource = "admin"
	// // cred.AuthMechanism = "SCRAM-SHA-256"
	// cred.Username = vars.DBUser
	// cred.Password = vars.DBPass

	uri := vars.DBConn
	opts := options.Client().ApplyURI(uri) //.SetAuth(cred)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal().Err(err).Str("uri", uri).Msg("failed to connect to local mongo server")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal().Err(err).Str("uri", uri).Msg("failed to ping mongo server")
	}

	log.Info().Str("uri", uri).Msg("connected to mongo server")

	return client, nil
}

func DisconnectMongo(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to disconnect from mongo server")
	}

	log.Info().Msg("disconnected from mongo server")
}
