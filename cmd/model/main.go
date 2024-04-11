package main

import (
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/recorders"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
)

var (
	app = kingpin.New("plot", "Plot market data along with model outputs.")

	dataDir   = app.Flag("datadir", "Data collection root dir").Required().String()
	graphFile = app.Flag("model", "Model JSON file path").Required().String()
	csvPath   = app.Flag("csv", "Path for a CSV output file").Required().String()
	start     = app.Flag("start", "start datetime formatted in RFC3339").String()
	end       = app.Flag("end", "end datetime formatted in RFC3339").String()
)

func main() {
	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	s, err := collection.NewDirStore(*dataDir)
	if err != nil {
		log.Fatal().Err(err).Str("dataDir", *dataDir).Msg("failed to make dir store")
	}

	c, err := collection.Load(s)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load collection")
	}

	g, err := blocks.LoadGraphFile(*graphFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load model file")
	}

	csvFile, err := os.Create(*csvPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to creat CSV file")
	}

	recorder := recorders.NewCSV(csvFile)

	if err = g.Init(recorder); err != nil {
		log.Fatal().Err(err).Msg("failed to init graph model")
	}

	ts := c.TimeSpan()
	tStart := ts.Start()
	tEnd := ts.End()

	if *start != "" {
		tStart, err = time.Parse(time.RFC3339, *start)
		if err != nil {
			log.Fatal().Err(err).Str("start", *start).Msg("failed to parse start time")
		}
	}

	if *end != "" {
		tEnd, err = time.Parse(time.RFC3339, *end)
		if err != nil {
			log.Fatal().Err(err).Str("end", *end).Msg("failed to parse end time")
		}
	}

	ts = timespan.NewTimeSpan(tStart, tEnd)

	bars, err := c.LoadBars(ts)
	if err != nil {
		log.Fatal().Err(err).
			Time("start", ts.Start()).
			Time("end", ts.End()).
			Msg("failed to load bars")
	}

	for _, bar := range bars {
		g.Update(bar)
	}

	recorder.Flush()
}
