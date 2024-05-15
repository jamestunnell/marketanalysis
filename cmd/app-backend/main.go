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
	DBName       = "marketanalysis"
	DefaultPort  = 4002
	DefaultDebug = false
)

type AppVariables struct {
	Debug  bool
	Port   int
	DBConn string
	// DBUser, DBPass string
}

type VarCandidate[T comparable] struct {
	Source string
	Value  T
}

func main() {
	vars := loadAppVars()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if vars.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	client := connectToLocalDB(vars)

	srv, router := server.New(vars.Port)

	loggingMiddleware := func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	}

	// router.Use(handlers.CORS(
	// 	handlers.AllowedOrigins([]string{"*"}),
	// 	handlers.AllowedMethods([]string{"*"}),
	// 	handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
	// ))
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
	// dbUser := app.Flag("dbuser", "Database user").Default("").String()
	// dbPass := app.Flag("dbpass", "Database password").Default("").String()

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	envvals, err := env.LoadValues()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load env values")
	}

	log.Info().Interface("values", envvals).Msg("loaded env values")

	vars.Port = loadAppVar[int]("port",
		newVarCandidate(*port, "CLI"),
		newVarCandidate(envvals.Port, "env"),
		newVarCandidate(DefaultPort, "default"))

	vars.Debug = loadAppVar[bool]("debug",
		newVarCandidate(*debug, "CLI"),
		newVarCandidate(envvals.Debug, "env"),
		newVarCandidate(DefaultDebug, "default"))

	vars.DBConn = loadAppVar[string]("dbconn",
		newVarCandidate(*dbConn, "CLI"),
		newVarCandidate(envvals.DBConn, "env"))

	// vars.DBUser = loadAppVar[string]("dbuser",
	// 	newVarCandidate(*dbUser, "CLI"),
	// 	newVarCandidate(envvals.DBUser, "env"))

	// vars.DBPass = loadAppVar[string]("dbpass",
	// 	newVarCandidate(*dbPass, "CLI"),
	// 	newVarCandidate(envvals.DBPass, "env"))

	return vars
}

func loadAppVar[T comparable](
	name string,
	first *VarCandidate[T],
	more ...*VarCandidate[T]) T {
	var zero T
	allSources := []string{}
	candidates := append([]*VarCandidate[T]{first}, more...)

	var source string
	var value T

	for _, c := range candidates {
		if c.Value != zero {
			value = c.Value
			source = c.Source

			break
		}

		allSources = append(allSources, c.Source)
	}

	if value == zero {
		log.Fatal().
			Str("name", name).
			Strs("sources", allSources).
			Msg("app var not found")
	}

	log.Info().
		Str("name", name).
		Str("source", source).
		Msgf("loaded app var")

	return value
}

func newVarCandidate[T comparable](val T, source string) *VarCandidate[T] {
	return &VarCandidate[T]{
		Source: source,
		Value:  val,
	}
}

func connectToLocalDB(vars *AppVariables) *mongo.Client {
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
