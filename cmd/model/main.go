package main

import (
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/recorders"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
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

	ts := c.GetTimeSpan()
	tStart := ts.Start()
	tEnd := ts.End()

	if *startStr != "" {
		startDate, err := date.Parse(date.RFC3339, *startStr)
		if err != nil {
			log.Fatal().Err(err).Str("start", *startStr).Msg("failed to parse start date")
		}

		tStart = startDate.In(loc)
	}

	if *endStr != "" {
		endDate, err := date.Parse(date.RFC3339, *endStr)
		if err != nil {
			log.Fatal().Err(err).Str("end", *endStr).Msg("failed to parse end date")
		}

		tEnd = endDate.In(loc)
	}

	ts = timespan.NewTimeSpan(tStart, tEnd)

	bars, err := c.LoadBars(ts)
	if err != nil {
		log.Fatal().Err(err).
			Time("start", ts.Start()).
			Time("end", ts.End()).
			Msg("failed to load bars")
	}

	log.Info().Err(err).
		Time("start", ts.Start()).
		Time("end", ts.End()).
		Int("count", len(bars)).
		Msg("loaded bars")

	for _, bar := range bars {
		g.Update(bar)
	}

	rec.Flush()
}
