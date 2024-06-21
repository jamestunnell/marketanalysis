package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
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
	"github.com/jamestunnell/marketanalysis/app/backend/background"
	"github.com/jamestunnell/marketanalysis/app/backend/env"
	"github.com/jamestunnell/marketanalysis/app/backend/server"
)

const (
	DBName       = "marketanalysis"
	DefaultPort  = "4002"
	DefaultDebug = "false"
)

type AppVariables struct {
	Debug  bool
	Port   int
	DBConn string
	// DBUser, DBPass string
}

type VarCandidate struct {
	Source string
	Value  string
}

func main() {
	vars := loadAppVars()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if vars.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	client := connectToLocalDB(vars)
	defer disconnectFromDB(client)

	db := client.Database(DBName)

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

	bg := background.NewSystem()

	api.BindAll(srv.GetRouter(), db, bg)

	srv.Start()
	defer srv.Stop()

	bg.Start()
	defer bg.Stop()

	server.BlockUntilSignaled(syscall.SIGINT, syscall.SIGTERM)
}

func loadAppVars() *AppVariables {
	vars := &AppVariables{}

	app := kingpin.New("backend server", "Provide market analysis features with an HTTP server`")
	debug := app.Flag("debug", "Enable debug mode").String()
	port := app.Flag("port", "Server port").String()
	dbConn := app.Flag("dbconn", "Database connection").String()
	// dbUser := backend.Flag("dbuser", "Database user").Default("").String()
	// dbPass := backend.Flag("dbpass", "Database password").Default("").String()

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	envvals, err := env.LoadValues()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load env values")
	}

	log.Info().Interface("values", envvals).Msg("loaded env values")

	vars.Port = loadAppVar[int](
		"port",
		strconv.Atoi,
		newVarCandidate(*port, "CLI"),
		newVarCandidate(os.Getenv(env.NamePort), "env"),
		newVarCandidate(DefaultPort, "default"))

	vars.Debug = loadAppVar[bool](
		"debug",
		strconv.ParseBool,
		newVarCandidate(*debug, "CLI"),
		newVarCandidate(os.Getenv(env.NameDebug), "env"),
		newVarCandidate(DefaultDebug, "default"))

	vars.DBConn = loadAppVar[string](
		"dbconn",
		func(s string) (string, error) { return s, nil },
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
	parse func(string) (T, error),
	first *VarCandidate,
	more ...*VarCandidate) T {
	allSources := []string{}
	candidates := append([]*VarCandidate{first}, more...)

	var source string
	var valStr string

	for _, c := range candidates {
		if c.Value != "" {
			valStr = c.Value
			source = c.Source

			break
		}

		allSources = append(allSources, c.Source)
	}

	if valStr == "" {
		log.Fatal().
			Str("name", name).
			Strs("sources", allSources).
			Msg("app var not found")
	}

	value, err := parse(valStr)
	if err != nil {
		log.Error().Err(err).Str("value", valStr).Msg("failed to parse app var")
	}

	log.Info().
		Str("name", name).
		Str("source", source).
		Interface("value", value).
		Msgf("loaded app var")

	return value
}

func newVarCandidate(val, source string) *VarCandidate {
	return &VarCandidate{
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
