package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/backtesting"
	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/strategies"
	"github.com/rs/zerolog/log"
)

type Command interface {
	Run() error
}

var (
	app       = kingpin.New("backtest", "Backtest trading strategy.")
	dataDir   = app.Flag("datadir", "Collection dir path.").Required().String()
	stratfile = app.Flag("stratfile", "Strategy JSON stratfile.").Required().String()
	verbose   = app.Flag("verbose", "Print all messages.").Bool()
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

	strat, err := strategies.LoadStrategyFromFile(*stratfile)
	if err != nil {
		log.Fatal().Err(err).Str("stratfile", *stratfile).Msg("failed to load strategy")
	}

	dc := backtesting.NewSequentialDates(c.Timespan().DateRangeIn(time.Local))
	dt := backtesting.NewDayTrader(c, strat, dc)
	positions := models.Positions{}

	fmt.Println("Running backtests")
	for dt.AnyLeft() {
		fmt.Printf(".")

		report, err := dt.RunTest()
		if err != nil {
			log.Fatal().Err(err).Msg("backtest failed")
		}

		positions = append(positions, report.Positions...)

		dt.Advance()
	}

	fmt.Printf("\n\n# Trades: %d\n", len(positions))

	analysis := positions.Analyze()

	fmt.Printf("Winning trades: %f%%\n", analysis.Winning*100.0)
	fmt.Printf("Total P/L: %f\n", analysis.TotalPL)
}
