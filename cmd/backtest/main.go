package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/strategies"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type Command interface {
	Run() error
}

const (
	DayTradeMarketOpenLocalMin  = 390
	DayTradeMarketCloseLocalMin = 780
)

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

	bars := c.GetBars(c.Timespan()).Local()

	start := bars[0].Timestamp.Round(24 * time.Hour)
	end := bars[len(bars)-1].Timestamp.Round(24 * time.Hour)

	totalPL := 0.0
	nWinning := 0
	nPositions := 0

	// test each day separately
	for dt := start; !dt.After(end); dt = dt.Add(24 * time.Hour) {
		yyyy, mm, dd := dt.Date()

		daytradeBars := sliceutils.Where(bars, func(b *models.Bar) bool {
			t := b.Timestamp
			dateMatch := (t.Year() == yyyy) && (t.Month() == mm) && (t.Day() == dd)
			min := 60*b.Timestamp.Hour() + t.Minute()
			timeMatch := min >= DayTradeMarketOpenLocalMin && min <= DayTradeMarketCloseLocalMin

			return dateMatch && timeMatch
		})

		slices.SortFunc(daytradeBars, func(a, b *models.Bar) bool {
			return a.Timestamp.Before(b.Timestamp)
		})

		if len(daytradeBars) > 0 {
			if *verbose {
				fmt.Printf("backtesting %d bars for %04d-%02d-%02d\n", len(daytradeBars), yyyy, mm, dd)
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

			backtest(strat, daytradeBars)

			positions := strat.ClosedPositions()

			if *verbose {
				fmt.Printf("Type\tEntry Time, Price\tExit Time, Price\tProfit/Loss\n")
			}
			for _, pos := range positions {
				pl, closed := pos.ClosedProfitLoss()
				if !closed {
					log.Fatal().Msg("position not closed")
				}

				if *verbose {
					fmt.Printf("%s\t%02d:%02d\t%f\t%02d:%02d\t%f\t%f\n",
						pos.Type(),
						pos.Entry().Time.Hour(), pos.Entry().Time.Minute(), pos.Entry().Price,
						pos.Exit().Time.Hour(), pos.Exit().Time.Minute(), pos.Exit().Price,
						pl)
				}

				totalPL += pl

				if pl > 0 {
					nWinning++
				}

				nPositions++
			}
		}
	}

	fmt.Printf("\n\nWinning: %d/%d (%f%%)\n", nWinning, nPositions, 100.0*float64(nWinning)/float64(nPositions))
	fmt.Printf("Total P/L: %f\n", totalPL)
}

func backtest(strat models.Strategy, bars []*models.Bar) {
	if strat.WarmupPeriod() > len(bars) {
		log.Info().Msg("not enough bars in day for warm-up and testing")

		return
	}

	wuBars := bars[:strat.WarmupPeriod()]
	remBars := bars[strat.WarmupPeriod():]

	if err := strat.WarmUp(wuBars); err != nil {
		log.Fatal().Err(err).Msg("failed to warm up strategy")
	}

	for _, bar := range remBars {
		strat.Update(bar)
	}

	strat.Close(bars[len(bars)-1])
}
