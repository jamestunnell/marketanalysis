package main

import (
	"encoding/json"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/aybabtme/uniplot/histogram"
	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/evaluation"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/processing"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type mlrData struct {
	Ins [][]float64
	Out []float64
}

const (
	DefaultHorizon      = 15
	DefaultTrainingDays = 90
)

var (
	app          = kingpin.New("backtest", "Backtest trading strategy.")
	debug        = app.Flag("debug", "Enable debug mode.").Bool()
	dataDir      = app.Flag("datadir", "Collection dir path.").Required().String()
	chainFile    = app.Flag("chainfile", "Processing chain JSON file.").Required().String()
	trainingDays = app.Flag("trainingDays", "Training days (Default is 90).").Int()
	inHor        = app.Flag("inhorizon", "Input horizon. Default is 15.").Int()
	outHor       = app.Flag("outhorizon", "Output horizon. Default is 15.").Int()
	futureHor    = app.Flag("futurehorizon", "Future horizon. Default is 15.").Int()

	seed = time.Now().Unix()
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if *trainingDays == 0 {
		*trainingDays = DefaultTrainingDays
	}

	if *inHor == 0 {
		*inHor = DefaultHorizon
	}

	if *outHor == 0 {
		*outHor = DefaultHorizon
	}

	if *futureHor == 0 {
		*futureHor = DefaultHorizon
	}

	s, err := collection.NewDirStore(*dataDir)
	if err != nil {
		log.Fatal().Err(err).Str("dataDir", *dataDir).Msg("failed to make dir store")
	}

	coll, err := collection.Load(s)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load collection")
	}

	chain := loadChain()
	randSource := rand.NewSource(seed)
	testingDays := int(math.Round(float64(coll.TimeSpan().Duration().Hours()/24.0) - float64(*trainingDays)))
	split := float64(*trainingDays) / (float64(coll.TimeSpan().Duration().Hours()) / 24.0)

	trainingBars, testingBars, err := evaluation.SplitCollectionRandomly(coll, split, randSource)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to split collection data")
	}

	pred := evaluation.NewSlopePredictor(*inHor, *outHor, *futureHor)

	err = pred.Train(chain, trainingBars)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to train predictor")
	}

	positions, err := pred.Test(chain, testingBars)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to test predictor")
	}

	profitLosses := sliceutils.Map(positions, func(pos *models.Position) float64 {
		return pos.ClosedPL
	})

	hist := histogram.Hist(20, profitLosses)
	_ = histogram.Fprint(os.Stdout, hist, histogram.Linear(30))

	q, err := stats.Quartile(profitLosses)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get quartiles")
	}

	analysis := positions.Analyze()

	log.Info().
		Float64("Q1", q.Q1).
		Float64("Q2", q.Q2).
		Float64("Q3", q.Q3).
		Int("training days", *trainingDays).
		Int("test days", testingDays).
		Float64("winning%", analysis.Winning*100.0).
		Float64("total profit/loss", analysis.TotalPL).
		Int("positions", len(positions)).
		Float64("avg position profit/loss", analysis.TotalPL/float64(len(positions))).
		Msg("evaluation complete")
}

func (d *mlrData) Inputs() [][]float64 {
	return d.Ins
}

func (d *mlrData) Output() []float64 {
	return d.Out
}

func loadChain() *processing.Chain {
	chainJSON, err := os.ReadFile(*chainFile)
	if err != nil {
		log.Fatal().Err(err).Str("fpath", *chainFile).Msg("failed to read chain JSON file")
	}

	var chain processing.Chain

	if err = json.Unmarshal(chainJSON, &chain); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal chain JSON")
	}

	log.Info().Str("json", string(chainJSON)).Msg("loaded chain")

	return &chain
}
