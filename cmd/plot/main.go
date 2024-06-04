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

	tzLoc     = backend.Flag("tz", `Timezone location`).Default("US/Pacific").String()
	dataDir   = backend.Flag("datadir", "Data collection root dir").Required().String()
	graphFile = backend.Flag("graphfile", "Graph Model JSON file path").Required().String()
	csvOut    = backend.Flag("csvout", "CSV output file").Required().String()
	start     = backend.Flag("start", "start datetime formatted in RFC3339").String()
	end       = backend.Flag("end", "end datetime formatted in RFC3339").String()
)

func main() {
	_ = kingpin.MustParse(backend.Parse(os.Args[1:]))

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

	csvFile, err := os.Create(*csvOut)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to creat CSV file")
	}

	loc, err := time.LoadLocation(*tzLoc)
	if err != nil {
		log.Fatal().Err(err).Str("tz", *tzLoc).Msg("failed to load timezone location")
	}

	recorder := recorders.NewCSV(csvFile, loc)

	if err = g.Init(recorder); err != nil {
		log.Fatal().Err(err).Msg("failed to init graph model")
	}

	ts := c.GetTimeSpan()
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
