package main

import (
	"fmt"
	"os"

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
	app     = kingpin.New("backtest", "Backtest trading strategy.")
	dataDir = app.Flag("datadir", "Collection dir path.").Required().String()
	verbose = app.Flag("verbose", "Print all messages.").Bool()

	tf     = app.Command("trendfollower", "Test a trend-following strategy.")
	tfFast = tf.Flag("fast", "Fast EMA period").Required().Int()
	tfSlow = tf.Flag("slow", "Slow EMA period").Required().Int()

	sc           = app.Command("scalper", "Test a scalping strategy.")
	scFast       = sc.Flag("fast", "Fast EMA period").Required().Int()
	scSlow       = sc.Flag("slow", "Slow EMA period").Required().Int()
	scTakeProfit = sc.Flag("takeprofit", "Take profit points").Required().Float64()
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
	case sc.FullCommand():
		params := models.Params{
			strategies.ParamFastPeriod: *scFast,
			strategies.ParamSlowPeriod: *scSlow,
			strategies.ParamTakeProfit: *scTakeProfit,
		}

		strat, err = strategies.NewScalper(params)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to make scalper strategy")
		}
	}

	dt := backtesting.NewDayTrader(c, strat)
	positions := []models.Position{}

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

	analysis := models.AnalyzePositions(positions)

	fmt.Printf("Winning trades: %f%%\n", analysis.Winning*100.0)
	fmt.Printf("Total P/L: %f\n", analysis.TotalPL)
}
