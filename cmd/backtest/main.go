package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/strategies"
	"github.com/rs/zerolog/log"
)

type Command interface {
	Run() error
}

var (
	app     = kingpin.New("backtest", "Backtest trading strategy.")
	dataDir = app.Flag("datadir", "Collection dir path.").Required().String()
	verbose = app.Flag("verbose", "Print all messages.").Bool()

	tf     = app.Command("trendfollower", "Test a trend-following strategy.")
	tfFast = tf.Flag("fast", "Fast EMA period").Required().Int()
	tfSlow = tf.Flag("slow", "Slow EMA period").Required().Int()
)

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	s, err := collection.NewDirStore(*dataDir)
	if err != nil {
		log.Fatal().Err(err).Str("dataDir", *dataDir).Msg("failed to make dir store")
	}

	c, err := collection.Load(s)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load collection")
	}

	var strat models.Strategy

	switch cmd {
	case tf.FullCommand():
		params := models.Params{
			strategies.ParamFastPeriod: *tfFast,
			strategies.ParamSlowPeriod: *tfSlow,
		}

		strat, err = strategies.NewTrendFollower(params)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to make trend follower strategy")
		}
	}

	bars := c.GetBars(c.Timespan())
	wuBars := bars[:strat.WarmupPeriod()]
	remBars := bars[strat.WarmupPeriod():]

	if err = strat.WarmUp(wuBars); err != nil {
		log.Fatal().Err(err).Msg("failed to warm up strategy")
	}

	for _, bar := range remBars {
		strat.Update(bar)
	}

	strat.Close(bars[len(bars)-1])

	positions := strat.ClosedPositions()
	totalPL := 0.0
	nWinning := 0

	if *verbose {
		fmt.Printf("Type\tEntry Time\tEntry Close\tExit Time\tExit Close\tProfit/Loss\n")
	}
	for _, pos := range positions {
		pl, closed := pos.ClosedProfitLoss()
		if !closed {
			log.Fatal().Msg("position not closed")
		}

		if *verbose {
			fmt.Printf("%s\t%v\t%f\t%v\t%f\t%f\n",
				pos.Type(),
				pos.Entry().Time.Local(), pos.Entry().Price,
				pos.Exit().Time.Local(), pos.Exit().Price,
				pl)
		}

		totalPL += pl

		if pl > 0 {
			nWinning++
		}
	}

	fmt.Printf("Winning: %d/%d (%f%%)\n", nWinning, len(positions), 100.0*float64(nWinning)/float64(len(positions)))
	fmt.Printf("Total P/L: %f\n", totalPL)
}
