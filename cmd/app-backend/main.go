package main

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/jamestunnell/marketanalysis/app/backend/api"
	"github.com/jamestunnell/marketanalysis/app/backend/env"
	"github.com/jamestunnell/marketanalysis/app/backend/server"
)

const (
	DBName      = "marketanalysis"
	DefaultPort = 4002
)

type AppVariables struct {
	Debug  bool
	Port   int
	DBConn string
}

func main() {
	vars := loadAppVars()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if vars.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	client := connectToLocalDB(vars.DBConn)

	srv, router := server.New(vars.Port)

	loggingMiddleware := func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	}

	router.Use(mux.MiddlewareFunc(loggingMiddleware))

	api.BindAll(srv.GetRouter(), client.Database(DBName))

	srv.Start()

	defer srv.Stop()
	defer disconnectFromDB(client)

	server.BlockUntilSignaled(syscall.SIGINT, syscall.SIGTERM)
}

func loadAppVars() *AppVariables {
	vars := &AppVariables{}

	app := kingpin.New("backend server", "Provide market analysis features with an HTTP server`")
	debug := app.Flag("debug", "Enable debug mode").Default("false").Bool()
	port := app.Flag("port", "Server port").Default("0").Int()
	dbConn := app.Flag("dbconn", "Database connection").Default("").String()

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	envvals, err := env.LoadValues()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load env values")
	}

	log.Info().Interface("values", envvals).Msg("loaded env values")

	var portSource string

	switch {
	case *port != 0:
		portSource = "CLI"

		vars.Port = *port
	case envvals.Port != 0:
		portSource = "env"

		vars.Port = envvals.Port
	default:
		portSource = "default"

		vars.Port = DefaultPort
	}

	var dbConnSource string

	switch {
	case *dbConn != "":
		dbConnSource = "CLI"

		vars.DBConn = *dbConn
	case envvals.DBConn != "":
		dbConnSource = "env"

		vars.DBConn = envvals.DBConn
	default:
		log.Fatal().Msg("dbconn not set through CLI or env")
	}

	var debugSource string

	switch {
	case *debug:
		debugSource = "CLI"

		vars.Debug = true
	case envvals.Debug:
		debugSource = "env"

		vars.Debug = true
	}

	log.Info().Int("port", vars.Port).Msgf("port assigned from %s value", portSource)
	log.Info().Str("dbconn", vars.DBConn).Msgf("dbconn assigned from %s value", dbConnSource)
	log.Info().Bool("debug", vars.Debug).Msgf("debug set in %s", debugSource)

	return vars
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
