package main

import (
	"net/http"
	"os"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app/backend/api"
	"github.com/jamestunnell/marketanalysis/app/backend/appvars"
	"github.com/jamestunnell/marketanalysis/app/backend/background"
	"github.com/jamestunnell/marketanalysis/app/backend/database"
	"github.com/jamestunnell/marketanalysis/app/backend/server"
)

func main() {
	vars, err := appvars.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load app vars")
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if vars.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	client, err := database.ConnectMongo(vars)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to mongo DB")
	}

	db := client.Database(database.MongoDBName)

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
