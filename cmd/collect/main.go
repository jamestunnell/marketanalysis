package main

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/commands/collectbars"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"
)

var (
	app = kingpin.New("collect", "Collect historical 1-minute bar data.`")

	startStr = app.Flag("start", "Start date formatted according to RFC3339.").Required().String()
	endStr   = app.Flag("end", "End date formatted according to RFC3339.").String()
	dir      = app.Flag("dir", "Collection dir path.").Required().String()
	sym      = app.Flag("sym", "The stock symbol.").Required().String()
)

func main() {
	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	startDate, err := date.Parse(date.RFC3339, *startStr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse start date")
	}

	startTime := startDate.UTC()

	var endDate date.Date

	if *endStr == "" {
		endDate = date.Today()
	} else {
		endDate, err = date.Parse(date.RFC3339, *endStr)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse end date")
		}
	}

	endTime := endDate.UTC()
	params := &collectbars.Params{
		Start:         startTime,
		End:           endTime,
		CollectionDir: *dir,
		Symbol:        *sym,
	}

	cmd, err := collectbars.New(params)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make collectbars command")
	}

	if err := cmd.Run(); err != nil {
		log.Fatal().Err(err).Msg("command failed")
	}
}
