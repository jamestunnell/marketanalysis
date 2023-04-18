package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/commands/collectbars"
	"github.com/rs/zerolog/log"
)

var (
	app = kingpin.New("collect", "Collect historical 1-minute bar data.`")

	start = app.Flag("start", "Start date-time formatted according to RFC3339.").Required().String()
	end   = app.Flag("end", "End date-time formatted according to RFC3339.").String()
	dir   = app.Flag("dir", "Collection dir path.").Required().String()
	sym   = app.Flag("sym", "The stock symbol.").Required().String()
)

func main() {
	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	tStart, err := time.Parse(time.RFC3339, *start)
	if err != nil {
		err = fmt.Errorf("failed to parse start datetime '%s' using RFC3339: %w", *start, err)

		log.Fatal().Err(err).Msg("failed to parse start time")
	}

	var tEnd time.Time

	if *end == "" {
		tEnd = time.Now().Add(-15 * time.Minute)
	} else {
		tEnd, err = time.Parse(time.RFC3339, *end)
		if err != nil {
			err = fmt.Errorf("failed to parse end datetime '%s' using RFC3339: %w", *end, err)

			log.Fatal().Err(err).Msg("failed to parse end time")
		}
	}

	params := &collectbars.Params{
		Start:         tStart,
		End:           tEnd,
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
