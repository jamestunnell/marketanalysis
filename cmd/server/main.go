package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/jamestunnell/marketanalysis/backend"
	"github.com/jamestunnell/marketanalysis/server"
)

const (
	DBName = "marketanalysis"
)

func main() {
	app := kingpin.New("server", "Provide market analysis features with an HTTP server`")
	debug := app.Flag("debug", "Enable debug mode").Default("false").Bool()
	port := app.Flag("port", "Server port").Required().Int()
	dbConn := app.Flag("dbconn", "Database connection").String()

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	client := connectToLocalDB(*dbConn)

	srv := server.NewServer(*port)

	backend.BindAPI(srv.GetRouter(), client.Database(DBName))

	srv.Start()

	defer srv.Stop()
	defer disconnectFromDB(client)

	server.BlockUntilSignaled(syscall.SIGINT, syscall.SIGTERM)
}

func connectToLocalDB(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal().Err(err).Str("uri", uri).Msg("failed to connect to local mongo server")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal().Err(err).Str("uri", uri).Msg("failed to ping mongo server")
	}

	log.Info().Str("uri", uri).Msg("connected to mongo server")

	return client
}

func disconnectFromDB(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to disconnect from mongo server")
	}

	log.Info().Msg("disconnected from mongo server")
}