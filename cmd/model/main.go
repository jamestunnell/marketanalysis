package main

import (
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/recorders"
	"github.com/rickb777/date"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	app = kingpin.New("plot", "Plot market data along with model outputs.")

	tz        = app.Flag("tz", `Timezone location`).Default("US/Pacific").String()
	debug     = app.Flag("debug", "Enable debug mode").Bool()
	dataDir   = app.Flag("datadir", "Data collection root dir").Required().String()
	modelFile = app.Flag("model", "Model JSON file path").Required().String()
	csvPath   = app.Flag("csv", "Path for a CSV output file").Required().String()
	startStr  = app.Flag("start", "start date (YYYY-MM-DD)").String()
	endStr    = app.Flag("end", "end date (YYYY-MM-DD)").String()
)

func main() {
	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	s, err := collection.NewDirStore(*dataDir)
	if err != nil {
		log.Fatal().Err(err).Str("dataDir", *dataDir).Msg("failed to make dir store")
	}

	c, err := collection.Load(s)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load collection")
	}

	g, err := blocks.LoadGraphFile(*modelFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load model file")
	}

	csvFile, err := os.Create(*csvPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to creat CSV file")
	}

	loc, err := time.LoadLocation(*tz)
	if err != nil {
		log.Fatal().Err(err).Str("tz", *tz).Msg("failed to load timezone location")
	}

	rec := recorders.NewCSV(csvFile, loc)

	if err = g.Init(rec); err != nil {
		log.Fatal().Err(err).Msg("failed to init graph model")
	}

	start := c.GetFirstDate()
	end := c.GetLastDate()

	if *startStr != "" {
		var err error

		start, err = date.Parse(date.RFC3339, *startStr)
		if err != nil {
			log.Fatal().Err(err).Str("start", *startStr).Msg("failed to parse start date")
		}
	}

	if *endStr != "" {
		var err error

		end, err = date.Parse(date.RFC3339, *endStr)
		if err != nil {
			log.Fatal().Err(err).Str("end", *endStr).Msg("failed to parse end date")
		}
	}

	bars, err := c.LoadBars(start, end)
	if err != nil {
		log.Fatal().Err(err).
			Stringer("start", start).
			Stringer("end", end).
			Msg("failed to load bars")
	}

	log.Info().Err(err).
		Stringer("start", start).
		Stringer("end", end).
		Int("count", len(bars)).
		Msg("loaded bars")

	for _, bar := range bars {
		g.Update(bar)
	}

	rec.Flush()
}
