package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/aybabtme/uniplot/histogram"
	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type (
	Params struct {
		Bins, Runs, Iter int
		Min, Max, Start  float64
	}

	Results struct {
		Min, Max, Mean, Median float64
	}

	processReturnFunc func(balance, ret float64) float64
)

const (
	BinsDefault = 20
	IterDefault = 20
	RunsDefault = 100000
)

var (
	app   = kingpin.New("sim-returns", "Simulate returns by repeated random experiment.")
	debug = app.Flag("debug", "Enable debug mode.").Bool()
	runs  = app.Flag("runs", "Number of runs. Default is 100,000.").Int()
	iter  = app.Flag("iter", "Number of iterations per run. Default is 20.").Int()
	bins  = app.Flag("bins", "Number of histogram bins. Default is 20.").Int()
	start = app.Flag("start", "Start balance").Required().Float64()
	min   = app.Flag("min", "Minimum ratio/loss amount").Required().Float64()
	max   = app.Flag("max", "Maximum ratio/profit amount").Required().Float64()

	add = app.Command("add", "Add profit/loss amount each iter.")
	mul = app.Command("mul", "Mul profit/loss ratio each iter.")
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	params := getParams()

	var processReturn processReturnFunc

	switch command {
	case add.FullCommand():
		processReturn = addProfitLoss
	case mul.FullCommand():
		processReturn = mulRatio
	}

	endBalances := make([]float64, params.Runs)

	for i := 0; i < params.Runs; i++ {
		endBalances[i] = simReturns(params, processReturn)
	}

	logStats(endBalances)
	printHist(endBalances, params.Bins)
}

func logStats(x []float64) {
	min, _ := stats.Min(x)
	max, _ := stats.Max(x)
	median, _ := stats.Median(x)
	mean, _ := stats.Mean(x)

	log.Info().
		Float64("min", min).
		Float64("max", max).
		Float64("median", median).
		Float64("mean", mean).
		Msg("stats")
}

func printHist(x []float64, bins int) {
	hist := histogram.Hist(bins, x)

	err := histogram.Fprint(os.Stdout, hist, histogram.Linear(5))
	if err != nil {
		log.Warn().Err(err).Msg("failed to print histogram")
	}
}

func getParams() *Params {
	if *bins == 0 {
		*bins = BinsDefault
	}

	if *runs == 0 {
		*runs = RunsDefault
	}

	if *iter == 0 {
		*iter = IterDefault
	}

	params := &Params{
		Runs:  *runs,
		Iter:  *iter,
		Start: *start,
		Min:   *min,
		Max:   *max,
		Bins:  *bins,
	}

	log.Info().
		Int("runs", params.Runs).
		Int("iter", params.Iter).
		Int("bins", params.Bins).
		Float64("min", params.Min).
		Float64("max", params.Max).
		Msg("params")

	return params
}

func simReturns(params *Params, processReturn processReturnFunc) float64 {
	balance := params.Start
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < params.Iter; i++ {
		ret := params.Min + (params.Max-params.Min)*rng.Float64()
		balance = processReturn(balance, ret)
	}

	return balance
}

func addProfitLoss(balance, profitLoss float64) float64 {
	return balance + profitLoss
}

func mulRatio(balance, ratio float64) float64 {
	return balance * ratio
}
