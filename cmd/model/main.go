package main

import (
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/commands"
	"github.com/jamestunnell/marketanalysis/commands/model"
	"github.com/jamestunnell/marketanalysis/recorders"
	"github.com/rickb777/date"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	app := kingpin.New("model", "Work with a model.")
	debug := backend.Flag("debug", "Enable debug mode").Bool()

	run := backend.Command("run", "Run the model with bar data from given date")

	runData := run.Flag("data", "Data collection root dir").Required().String()
	runModel := run.Flag("model", "Model JSON file path").Required().String()
	runCSV := run.Flag("csv", "Path for a CSV output file").Required().String()
	runDate := run.Flag("date", "Date to test (YYYY-MM-DD)").String()
	runTZ := run.Flag("tz", "Recording time zone location").String()

	cmdName := kingpin.MustParse(backend.Parse(os.Args[1:]))

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	var cmd commands.Command

	switch cmdName {
	case run.FullCommand():
		d, err := date.Parse(date.RFC3339, *runDate)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse date")
		}

		c, err := collection.LoadFromDir(*runData)
		if err != nil {
			log.Fatal().Err(err).Str("dir", *runData).Msg("failed to load collection")
		}

		csvFile, err := os.Create(*runCSV)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to creat CSV file")
		}

		// use collection location by default
		recordLoc := c.GetLocation()

		if *runTZ != "" {
			loc, err := time.LoadLocation(*runTZ)
			if err != nil {
				log.Fatal().Err(err).Str("tz", *runTZ).Msg("failed to load recording timezone location")
			}

			recordLoc = loc
		}

		rec := recorders.NewCSV(csvFile, recordLoc)

		cmd = &model.ModelRun{
			Collection: c,
			ModelFile:  *runModel,
			Recorder:   rec,
			Date:       d,
		}
	default:
		log.Fatal().Msgf("unknown command %s", cmdName)
	}

	commands.InitAndRun(cmd)
}
