package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"

	"github.com/MaxHalford/eaopt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/backtesting"
	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/optimization"
	"github.com/jamestunnell/marketanalysis/params"
	"github.com/jamestunnell/marketanalysis/strategies"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
)

type PartsToStrategyFunc func(parts []optimization.PartialGenome) (models.Strategy, error)

const (
	DefaultMaxPeriod = 200
	DefaultTestDays  = 100

	DefaultMinTakeProfit = 0.5
	DefaultMaxTakeProfit = 2.0
)

var (
	app       = kingpin.New("optimize", "Optimize trading strategy.")
	dataDir   = app.Flag("datadir", "Collection dir path.").Required().String()
	maxPeriod = app.Flag("maxperiod", "Maximum MA period. Default is 200.").Int()
	testDays  = app.Flag("testdays", "Number of random days to use for testing. Default is 100.").Int()
	outFile   = app.Flag("outfile", "Filepath for stategy JSON file. Default is strategy.json").String()

	tf = app.Command("trendfollower", "Optimize a trend-following strategy.")

	sc            = app.Command("scalper", "Optimize a scalper strategy.")
	minTakeProfit = sc.Flag("mintakeprofit", "Maximum take profit param value. Default is 0.5").Float()
	maxTakeProfit = sc.Flag("maxtakeprofit", "Maximum take profit param value. Default is 2.0").Float()
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

	// strat, err := strategies.LoadStrategyFromFile(*stratFile)
	// if err != nil {
	// 	log.Fatal().Err(err).Str("file", *stratFile).Msg("failed to load strategy from file")
	// }

	if *maxPeriod == 0 {
		*maxPeriod = DefaultMaxPeriod
	}

	if *testDays == 0 {
		*testDays = DefaultTestDays
	}

	if *outFile == "" {
		*outFile = "./strategy.json"
	}

	var newGenome optimization.NewGenomeFunc

	var partsToStrategy PartsToStrategyFunc

	switch cmd {
	case tf.FullCommand():
		partsToStrategy = func(parts []optimization.PartialGenome) (models.Strategy, error) {
			periods := parts[0].(*optimization.OrderedInts)
			params := models.Params{
				strategies.ParamFastPeriod: params.NewInt(periods.Ints[0]),
				strategies.ParamSlowPeriod: params.NewInt(periods.Ints[1]),
			}

			s, err := strategies.NewTrendFollower(params)
			if err != nil {
				return nil, fmt.Errorf("failed to make trend follower strategy: %w", err)
			}

			return s, nil
		}
		newGenome = func(rng *rand.Rand) eaopt.Genome {
			periods := optimization.RandomOrderedInts(2, rng, 1, *maxPeriod)

			return optimization.NewCompositeGenome(makeFit(partsToStrategy, c), periods)
		}
	case sc.FullCommand():
		if *minTakeProfit == 0.0 {
			*minTakeProfit = DefaultMinTakeProfit
		}

		if *maxTakeProfit == 0.0 {
			*maxTakeProfit = DefaultMaxTakeProfit
		}

		partsToStrategy = func(parts []optimization.PartialGenome) (models.Strategy, error) {
			periods := parts[0].(*optimization.OrderedInts)
			takeProfit := parts[1].(*optimization.FloatValue)

			params := models.Params{
				strategies.ParamFastPeriod: params.NewInt(periods.Ints[0]),
				strategies.ParamSlowPeriod: params.NewInt(periods.Ints[1]),
				strategies.ParamTakeProfit: params.NewFloat(takeProfit.Value),
			}

			s, err := strategies.NewScalper(params)
			if err != nil {
				return nil, fmt.Errorf("failed to make scalper strategy: %w", err)
			}

			return s, nil
		}
		newGenome = func(rng *rand.Rand) eaopt.Genome {
			periods := optimization.RandomOrderedInts(2, rng, 1, *maxPeriod)
			takeProfit := optimization.RandomFloatValue(rng, *minTakeProfit, *maxTakeProfit)

			return optimization.NewCompositeGenome(makeFit(partsToStrategy, c), periods, takeProfit)
		}
	}

	if newGenome == nil {
		log.Fatal().Msg("unsupported stategy type")
	}

	config := eaopt.NewDefaultGAConfig()

	config.HofSize = 5
	config.PopSize = 100
	config.NGenerations = 300
	config.ParallelEval = true

	best, err := optimization.EAOpt(newGenome, config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to optimize")
	}

	fmt.Println("\nHall of Fame:")
	fmt.Println("Fitness\tStrategy JSON")
	for i := 0; (i < 10) && (i < len(best)); i++ {
		cg := best[i].Genome.(*optimization.CompositeGenome)
		s, err := partsToStrategy(cg.Parts)
		if err != nil {
			log.Warn().Err(err).Msg("failed to make strategy with composite genome parts")
		}

		var buf bytes.Buffer

		strategies.StoreStrategy(s, &buf)

		fmt.Printf("%f\t%#v\n", best[i].Fitness, buf.String())

		if i == 0 {
			err = strategies.StoreStrategyToFile(s, *outFile)
			if err != nil {
				log.Warn().Err(err).Msg("failed to store best strategy")
			}
		}
	}

}

func makeFit(f PartsToStrategyFunc, c collection.Collection) optimization.CompositeFitFunc {
	const seed int64 = 100553122

	return func(parts []optimization.PartialGenome) (float64, error) {
		ts := c.Timespan()
		dateRange := timespan.NewDateRangeOf(
			ts.Start(), ts.End().Sub(ts.Start()))

		s, err := f(parts)
		if err != nil {
			return 0.0, fmt.Errorf("failed to make strategy: %w", err)
		}

		randomDates := backtesting.NewRandomDates(
			dateRange, *testDays, rand.NewSource(seed))
		trader := backtesting.NewDayTrader(c, s, randomDates)

		fitVal, err := evaluate(trader)
		if err != nil {
			err = fmt.Errorf("failed to evaluate: %w", err)

			return 0.0, err
		}

		var buf bytes.Buffer

		strategies.StoreStrategy(s, &buf)

		log.Debug().
			// Str("strategy", buf.String()).
			Float64("fitness", fitVal).
			Msg("evaluated strategy")

		return fitVal, nil
	}
}

func evaluate(tester backtesting.Tester) (float64, error) {
	positions := models.Positions{}

	for tester.AnyLeft() {

		report, err := tester.RunTest()
		if err != nil {
			return 0.0, fmt.Errorf("backtest failed: %w", err)
		}

		positions = append(positions, report.Positions...)

		tester.Advance()
	}

	analysis := positions.Analyze()
	fitness := 1.0 - analysis.Winning

	return fitness, nil
}
