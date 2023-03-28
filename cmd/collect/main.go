package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/commands/collectbars"
	"github.com/rs/zerolog/log"
)

type Command interface {
	Run() error
}

var (
	app = kingpin.New("collect", "Collect market data.")

	bars = app.Command("bars", "Collect historical 1-minute bar data.")

	barsStart = bars.Flag("start", "Start date-time formatted according to RFC3339.").Required().String()
	barsEnd   = bars.Flag("end", "End date-time formatted according to RFC3339.").Required().String()
	barsDir   = bars.Flag("dir", "Collection dir path.").Required().String()
	barsSym   = bars.Flag("sym", "The stock symbol.").Required().String()
)

func main() {
	var cmd Command

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case bars.FullCommand():
		tStart, err := time.Parse(time.RFC3339, *barsStart)
		if err != nil {
			err = fmt.Errorf("failed to parse start datetime '%s' using RFC3339: %w", *barsStart, err)

			log.Fatal().Err(err).Msg("failed to parse start time")
		}

		tEnd, err := time.Parse(time.RFC3339, *barsEnd)
		if err != nil {
			err = fmt.Errorf("failed to parse end datetime '%s' using RFC3339: %w", *barsEnd, err)

			log.Fatal().Err(err).Msg("failed to parse end time")
		}

		params := &collectbars.Params{
			Start:         tStart,
			End:           tEnd,
			CollectionDir: *barsDir,
			Symbol:        *barsSym,
		}
		cmd, err = collectbars.New(params)

		if err != nil {
			log.Fatal().Err(err).Msg("failed to make collectbars command")
		}
	}

	if cmd == nil {
		log.Fatal().Msg("unknown command")
	}

	if err := cmd.Run(); err != nil {
		log.Fatal().Err(err).Msg("command failed")
	}
}

func parseStartEndTimes(start, end string) (tStart, tEnd time.Time, err error) {
	tStart, err = time.Parse(time.RFC3339, start)
	if err != nil {
		err = fmt.Errorf("failed to parse start datetime '%s' using RFC3339: %w", start, err)

		return
	}

	tEnd, err = time.Parse(time.RFC3339, end)
	if err != nil {
		err = fmt.Errorf("failed to parse end datetime '%s' using RFC3339: %w", end, err)

		return
	}

	return
}
