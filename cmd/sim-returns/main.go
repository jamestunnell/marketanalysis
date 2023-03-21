package main

import (
	"fmt"
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

	results := make([]float64, params.Runs)

	for i := 0; i < params.Runs; i++ {
		results[i] = simReturns(params, processReturn)
	}

	printStats(results)
	printProbs(results, params.Start)
	printHist(results, params.Bins)
}

func printStats(x []float64) {
	min, _ := stats.Min(x)
	max, _ := stats.Max(x)
	median, _ := stats.Median(x)
	mean, _ := stats.Mean(x)
	stddev, _ := stats.StandardDeviation(x)

	fmt.Println("Stats:")
	fmt.Printf("  Min: %f\n", min)
	fmt.Printf("  Max: %f\n", max)
	fmt.Printf("  Median: %f\n", median)
	fmt.Printf("  Mean: %f\n", mean)
	fmt.Printf("  Std Dev: %f\n", stddev)
	fmt.Println("")
}

func printProbs(vals []float64, start float64) {
	fmt.Println("Loss Probabilities:")
	fmt.Printf("  25%% or more: %f%%\n", 100.0*probLE(vals, start*0.75))
	fmt.Printf("  50%% or more: %f%%\n", 100.0*probLE(vals, start*0.5))
	fmt.Printf("  75%% or more: %f%%\n", 100.0*probLE(vals, start*0.25))
	fmt.Printf("  100%% or more: %f%%\n", 100.0*probLE(vals, 0))
	fmt.Println("")

	fmt.Println("Profit Probabilities:")
	fmt.Printf("  25%% or more: %f%%\n", 100.0*probGE(vals, start*1.25))
	fmt.Printf("  50%% or more: %f%%\n", 100.0*probGE(vals, start*1.5))
	fmt.Printf("  75%% or more: %f%%\n", 100.0*probGE(vals, start*1.75))
	fmt.Printf("  100%% or more: %f%%\n", 100.0*probGE(vals, start*2))
	fmt.Println("")
}

func probLE(vals []float64, level float64) float64 {
	le := func(val float64) bool {
		return val <= level
	}

	return float64(count(vals, le)) / float64(len(vals))
}

func probGE(vals []float64, level float64) float64 {
	ge := func(val float64) bool {
		return val >= level
	}

	return float64(count(vals, ge)) / float64(len(vals))
}

func count[T any](ts []T, f func(t T) bool) int {
	count := 0

	for _, t := range ts {
		if f(t) {
			count += 1
		}
	}

	return count
}

func printHist(x []float64, bins int) {
	fmt.Println("Histogram:")
	hist := histogram.Hist(bins, x)

	err := histogram.Fprint(os.Stdout, hist, histogram.Linear(5))
	if err != nil {
		log.Warn().Err(err).Msg("failed to print histogram")
	}

	fmt.Println("")
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

	fmt.Println("Params:")
	fmt.Printf("  Runs: %d\n", params.Runs)
	fmt.Printf("  Iter: %d\n", params.Iter)
	fmt.Printf("  Bins: %d\n", params.Bins)
	fmt.Printf("  Min: %f\n", params.Min)
	fmt.Printf("  Max: %f\n", params.Max)
	fmt.Println("")

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
