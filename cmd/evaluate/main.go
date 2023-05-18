package main

import (
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/evaluation"
	"github.com/jamestunnell/marketanalysis/processing"
	"github.com/jamestunnell/marketanalysis/regression/linregression"
	"github.com/jamestunnell/marketanalysis/regression/mlregression"
	"github.com/jamestunnell/marketanalysis/util/dateutils"
)

type mlrData struct {
	Ins [][]float64
	Out []float64
}

const (
	DayTradeMarketOpenLocalMin  = 390
	DayTradeMarketCloseLocalMin = 780

	DefaultSplit = 0.5
)

var (
	app       = kingpin.New("backtest", "Backtest trading strategy.")
	debug     = app.Flag("debug", "Enable debug mode.").Bool()
	dataDir   = app.Flag("datadir", "Collection dir path.").Required().String()
	chainFile = app.Flag("chainfile", "Processing chain JSON file.").Required().String()
	split     = app.Flag("split", "Training/testing split (0.0 to 1.0). Default is 0.5").Float64()

	seed = time.Now().Unix()
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if *split == 0.0 {
		*split = DefaultSplit
	}

	s, err := collection.NewDirStore(*dataDir)
	if err != nil {
		log.Fatal().Err(err).Str("dataDir", *dataDir).Msg("failed to make dir store")
	}

	coll, err := collection.Load(s)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load collection")
	}

	chain, err := processing.LoadChainFromFile(*chainFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load predictor")
	}

	log.Info().Msg("loaded chain")

	randSource := rand.NewSource(seed)
	dateRange := coll.TimeSpan().DateRangeIn(time.Local)
	numTrain := int(*split * float64(dateRange.Days()))
	dateCtrl := dateutils.NewRandomDates(dateRange, numTrain, randSource)
	inSlopes := []float64{}
	inIntercepts := []float64{}
	outSlopes := []float64{}

	for dateCtrl.AnyLeft() {
		open := dateCtrl.Current().Local().Add(DayTradeMarketOpenLocalMin * time.Minute)
		close := dateCtrl.Current().Local().Add(DayTradeMarketCloseLocalMin * time.Minute)
		ts := timespan.NewTimeSpan(open, close)

		bars, err := coll.LoadBars(ts)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load bars")
		}

		if len(bars) != 0 {
			log.Debug().
				Str("date", dateCtrl.Current().Format(date.RFC3339)).
				Int("num bars", len(bars)).
				Msg("training")

			err = evaluation.ExtractSlopes(
				chain, bars, 15, 15, func(inLine, outLine *linregression.Line) {
					inSlopes = append(inSlopes, inLine.Slope)
					inIntercepts = append(inIntercepts, inLine.Intercept)
					outSlopes = append(outSlopes, outLine.Slope)
				})

			if err != nil {
				log.Fatal().Err(err).Msg("failed to train on bars")
			}
		}

		dateCtrl.Advance()
	}

	l := mlregression.NewSliceLearner()
	d := &mlrData{Ins: [][]float64{inSlopes, inIntercepts}, Out: outSlopes}

	pred, err := l.Learn(d, 0.01, 100)
	if err != nil {
		log.Fatal().Err(err).Msg("ML regression failed")
	}

	dateCtrl = dateutils.NewRandomDates(dateRange, int(dateRange.Days())-numTrain, randSource)

	sumAbsErrs := 0.0
	numAbsErrs := 0

	for dateCtrl.AnyLeft() {
		open := dateCtrl.Current().Local().Add(DayTradeMarketOpenLocalMin * time.Minute)
		close := dateCtrl.Current().Local().Add(DayTradeMarketCloseLocalMin * time.Minute)
		ts := timespan.NewTimeSpan(open, close)

		bars, err := coll.LoadBars(ts)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load bars")
		}

		if len(bars) != 0 {
			log.Debug().
				Str("date", dateCtrl.Current().Format(date.RFC3339)).
				Int("num bars", len(bars)).
				Msg("testing")

			err = evaluation.ExtractSlopes(
				chain, bars, 15, 15, func(inLine, outLine *linregression.Line) {
					predIns := []float64{inLine.Slope, inLine.Intercept}
					predOut := pred.PredictOne(predIns)

					sumAbsErrs += math.Abs(outLine.Slope - predOut)
					numAbsErrs++
				})

			if err != nil {
				log.Fatal().Err(err).Msg("failed to train on bars")
			}
		}

		dateCtrl.Advance()
	}

	avgAbsErr := sumAbsErrs / float64(numAbsErrs)

	log.Info().Float64("avg abs err", avgAbsErr).Msg("evaluation complete")
}

func (d *mlrData) Inputs() [][]float64 {
	return d.Ins
}

func (d *mlrData) Output() []float64 {
	return d.Out
}
