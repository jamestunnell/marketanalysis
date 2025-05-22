package appvars

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/app/backend/env"
	"github.com/rs/zerolog/log"
)

func Load() (*Vars, error) {
	debugStr := kingpin.Flag("debug", "Enable debug mode").String()
	portStr := kingpin.Flag("port", "Server port").String()
	dbConnStr := kingpin.Flag("dbconn", "Database connection").String()
	// dbUser := backend.Flag("dbuser", "Database user").Default("").String()
	// dbPass := backend.Flag("dbpass", "Database password").Default("").String()

	_ = kingpin.Parse()

	envvals, err := env.LoadValues()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load env values")
	}

	log.Info().Interface("values", envvals).Msg("loaded env values")

	port, err := loadAppVar[int](
		"port",
		strconv.Atoi,
		newVarCandidate(*portStr, "CLI"),
		newVarCandidate(os.Getenv(env.NamePort), "env"),
		newVarCandidate(DefaultPort, "default"))
	if err != nil {
		return nil, fmt.Errorf("failed to load port var: %w", err)
	}

	debug, err := loadAppVar[bool](
		"debug",
		strconv.ParseBool,
		newVarCandidate(*debugStr, "CLI"),
		newVarCandidate(os.Getenv(env.NameDebug), "env"),
		newVarCandidate(DefaultDebug, "default"))
	if err != nil {
		return nil, fmt.Errorf("failed to load debug var: %w", err)
	}

	dbConn, err := loadAppVar[string](
		"dbconn",
		func(s string) (string, error) { return s, nil },
		newVarCandidate(*dbConnStr, "CLI"),
		newVarCandidate(envvals.DBConn, "env"))
	if err != nil {
		return nil, fmt.Errorf("failed to load dbConn var: %w", err)
	}

	vars := &Vars{
		Port:   port,
		Debug:  debug,
		DBConn: dbConn,
	}

	return vars, nil
}
