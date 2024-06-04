package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"math"
// 	"math/rand"
// 	"os"
// 	"time"

// 	"github.com/MaxHalford/eaopt"
// 	"github.com/alecthomas/kingpin/v2"
// 	"github.com/jamestunnell/marketanalysis/backtesting"
// 	"github.com/jamestunnell/marketanalysis/collection"
// 	"github.com/jamestunnell/marketanalysis/models"
// 	"github.com/jamestunnell/marketanalysis/optimization"
// 	"github.com/jamestunnell/marketanalysis/predictors"
// 	"github.com/jamestunnell/marketanalysis/util/sliceutils"
// 	"github.com/rickb777/date/timespan"
// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"
// )

// const (
// 	DefaultTestDays = 100
// 	DefaultNumGens  = 50
// 	DefaultPopSize  = 15
// 	DefaultNumPops  = 2
// )

// var (
// 	app         = kingpin.New("optimize", "Optimize trading strategy.")
// 	debug       = backend.Flag("debug", "Enable debug mode.").Bool()
// 	dataDir     = backend.Flag("datadir", "Collection dir path.").Required().String()
// 	testDays    = backend.Flag("testdays", "Number of random days to use for testing. Default is 100.").Int()
// 	outFile     = backend.Flag("outfile", "Filepath for predictor JSON file. Default is ./predictor.json").String()
// 	upperLimits = backend.Flag("upperlimits", "JSON snippet to override default upper limits for params.").String()
// 	predName    = backend.Flag("predictor", "Predictor type name").Required().String()
// 	numGens     = backend.Flag("numgens", "Num generations for EA optimization. Default is 50.").Uint()
// 	popSize     = backend.Flag("popsize", "Population size for EA optimization. Default is 15.").Uint()
// 	numPops     = backend.Flag("numpops", "Num populations for EA optimization. Default is 2.").Uint()
// 	parallel    = backend.Flag("parallel", "Run EA optimization in parallel. Default is false.").Bool()

// 	seed = time.Now().Unix()
// )

// func main() {
// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

// 	zerolog.SetGlobalLevel(zerolog.InfoLevel)

// 	_ = kingpin.MustParse(backend.Parse(os.Args[1:]))

// 	if *debug {
// 		zerolog.SetGlobalLevel(zerolog.DebugLevel)
// 	}

// 	s, err := collection.NewDirStore(*dataDir)
// 	if err != nil {
// 		log.Fatal().Err(err).Str("dataDir", *dataDir).Msg("failed to make dir store")
// 	}

// 	c, err := collection.Load(s)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to load collection")
// 	}

// 	// strat, err := strategies.LoadStrategyFromFile(*stratFile)
// 	// if err != nil {
// 	// 	log.Fatal().Err(err).Str("file", *stratFile).Msg("failed to load strategy from file")
// 	// }

// 	if *testDays == 0 {
// 		*testDays = DefaultTestDays
// 	}

// 	if *outFile == "" {
// 		*outFile = "./strategy.json"
// 	}

// 	if *numGens == 0 {
// 		*numGens = DefaultNumGens
// 	}

// 	if *popSize == 0 {
// 		*popSize = DefaultPopSize
// 	}

// 	if *numPops == 0 {
// 		*numPops = DefaultNumPops
// 	}

// 	if *upperLimits != "" {
// 		overrideUpperLimits(*upperLimits)
// 	}

// 	newPredictor, found := predictors.GetNewPredictorFunc(*predName)
// 	if !found {
// 		log.Fatal().Msg("unknown predictor type")
// 	}

// 	gToP, newGenome := makeEAFuncs(newPredictor, c)

// 	config := eaopt.NewDefaultGAConfig()

// 	config.HofSize = 5
// 	config.PopSize = 100
// 	config.NGenerations = 300
// 	config.ParallelEval = true

// 	best, err := optimization.EAOpt(newGenome, gToP, config)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to optimize")
// 	}

// 	fmt.Println("\nHall of Fame:")
// 	fmt.Println("Fitness\tStrategy JSON")
// 	for i := 0; (i < 10) && (i < len(best)); i++ {
// 		g := best[i].Genome

// 		pred, err := gToP(g)
// 		if err != nil {
// 			log.Warn().
// 				Err(err).
// 				Interface("genome", g).
// 				Msg("failed to make predictor with genome")
// 		}

// 		var buf bytes.Buffer

// 		predictors.StorePredictor(pred, &buf)

// 		fmt.Printf("%f\t%s\n", best[i].Fitness, buf.String())

// 		if i == 0 {
// 			err = predictors.StorePredictorToFile(pred, *outFile)
// 			if err != nil {
// 				log.Warn().Err(err).Msg("failed to store best predictor")
// 			}
// 		}
// 	}

// }

// func overrideUpperLimits(jsonStr string) {
// 	var overrides map[string]any
// 	if err := json.Unmarshal([]byte(jsonStr), &overrides); err != nil {
// 		log.Fatal().
// 			Err(err).
// 			Msg("failed to unmarshal upper limits JSON")
// 	}

// 	for paramName, val := range overrides {
// 		lim, found := predictors.GetUpperLimit(paramName)
// 		if !found {
// 			log.Warn().
// 				Str("param", paramName).
// 				Msg("no upper limit found for param")
// 		}

// 		if !lim.Change(val) {
// 			log.Warn().
// 				Interface("value", val).
// 				Str("param", paramName).
// 				Msg("failed to change upper limit")
// 		}
// 	}
// }

// func makeEAFuncs(
// 	newPred predictors.NewPredictorFunc,
// 	c models.Collection,
// ) (optimization.GenomeToPredictorFunc, optimization.NewGenomeFunc) {
// 	gToP := func(g eaopt.Genome) (models.Predictor, error) {
// 		cg := g.(*optimization.CompositeGenome)
// 		pred := newPred()
// 		params := pred.Params()
// 		names := params.SortedNames()

// 		for i, name := range names {
// 			err := params[name].SetVal(cg.Parts[i].Value())
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to set param '%s': %w", name, err)
// 			}
// 		}

// 		return pred, nil
// 	}

// 	fit := makeFit(gToP, c)
// 	examplePred := newPred()

// 	newG := func(rng *rand.Rand) eaopt.Genome {
// 		params := examplePred.Params()
// 		names := params.SortedNames()
// 		parts := make([]optimization.PartialGenome, len(names))

// 		for i, name := range names {
// 			partial, err := optimization.MakePartial(params[name], rng)
// 			if err != nil {
// 				log.Fatal().
// 					Err(err).
// 					Str("param", name).
// 					Msg("failed to make param partial")
// 			}

// 			parts[i] = partial
// 		}

// 		return optimization.NewCompositeGenome(fit, parts...)
// 	}

// 	return gToP, newG
// }

// func makeFit(
// 	gToP optimization.GenomeToPredictorFunc,
// 	c models.Collection) optimization.FitFunc {
// 	return func(g eaopt.Genome) (float64, error) {
// 		ts := c.GetTimeSpan()
// 		dateRange := timespan.NewDateRangeOf(
// 			ts.Start(), ts.End().Sub(ts.Start()))

// 		pred, err := gToP(g)
// 		if err != nil {
// 			return 0.0, fmt.Errorf("failed to make predictor: %w", err)
// 		}

// 		sum := 0.0
// 		totalBars := 0
// 		evalPredictions := func(dir models.Direction, bars models.Bars) {
// 			totalBars += len(bars)
// 			sum += evalPrediction(dir, bars)
// 		}

// 		randomDates := backtesting.NewRandomDates(
// 			dateRange, *testDays, rand.NewSource(seed))
// 		trader := backtesting.NewDayTrader(c, pred, randomDates, evalPredictions)

// 		for trader.AnyLeft() {
// 			if err := trader.RunTest(); err != nil {
// 				return 0.0, fmt.Errorf("backtest failed: %w", err)
// 			}

// 			trader.Advance()
// 		}

// 		fitVal := -sum / float64(totalBars)

// 		log.Trace().
// 			Str("params", pred.Params().String()).
// 			Float64("fitness", fitVal).
// 			Msg("evaluated strategy")

// 		return fitVal, nil
// 	}
// }

// func evalPrediction(dir models.Direction, bars models.Bars) float64 {
// 	if len(bars) < 22 {
// 		log.Fatal().
// 			Int("len", len(bars)).
// 			Msg("not enough bars to evaluate prediction")
// 	}

// 	entry := bars[0].Close
// 	diffs := []float64{}
// 	remBars := bars[1:]

// 	switch dir {
// 	case models.DirUp:
// 		diffs = sliceutils.Map(remBars, func(b *models.Bar) float64 {
// 			return b.Close - entry
// 		})
// 	case models.DirDown:
// 		diffs = sliceutils.Map(remBars, func(b *models.Bar) float64 {
// 			return entry - b.Close
// 		})
// 	case models.DirNone:
// 		diffs = sliceutils.Map(remBars, func(b *models.Bar) float64 {
// 			return -math.Abs(entry - b.Close)
// 		})
// 	}

// 	return sumFloats(diffs)
// }

// func sumFloats(vals []float64) float64 {
// 	sum := 0.0

// 	for _, val := range vals {
// 		sum += val
// 	}

// 	return sum
// }
