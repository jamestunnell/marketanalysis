package main

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/rickb777/date"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/commands"
	"github.com/jamestunnell/marketanalysis/commands/collect"
)

func main() {
	app := kingpin.New("collect", "Collect historical 1-minute bar data.`")
	debug := backend.Flag("debug", "Enable debug mode").Bool()

	new := backend.Command("new", "Start a new collection")
	newDir := new.Flag("dir", "Collection dir path (created if needed).").Required().String()
	newSym := new.Flag("sym", "The stock symbol.").Required().String()
	newTZ := new.Flag("tz", "Time zone location").Default("America/New_York").String()
	newStart := new.Flag("start", "Start date formatted according to RFC3339.").Required().String()

	update := backend.Command("update", "Update existing collection data with the latest bars")
	updateDir := update.Flag("dir", "Existing collection dir path.").Required().String()

	cmdName := kingpin.MustParse(backend.Parse(os.Args[1:]))

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	var cmd commands.Command

	switch cmdName {
	case new.FullCommand():
		startDate, err := date.Parse(date.RFC3339, *newStart)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse start date")
		}

		if *newSym == "" {
			log.Fatal().Msg("no symbol given")
		}

		cmd = &collect.CollectNew{
			StartDate: startDate,
			Dir:       *newDir,
			Symbol:    *newSym,
			TimeZone:  *newTZ,
		}
	case update.FullCommand():
		cmd = &collect.CollectUpdate{
			Dir: *updateDir,
		}
	default:
		log.Fatal().Msgf("unknown command %s", cmdName)
	}

	commands.InitAndRun(cmd)
}
