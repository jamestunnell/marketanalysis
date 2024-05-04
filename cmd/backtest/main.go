package main

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/alecthomas/kingpin/v2"
// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"

// 	"github.com/jamestunnell/marketanalysis/backtesting"
// 	"github.com/jamestunnell/marketanalysis/collection"
// 	"github.com/jamestunnell/marketanalysis/models"
// 	"github.com/jamestunnell/marketanalysis/predictors"
// )

// var (
// 	app           = kingpin.New("backtest", "Backtest trading strategy.")
// 	debug         = app.Flag("debug", "Enable debug mode.").Bool()
// 	dataDir       = app.Flag("datadir", "Collection dir path.").Required().String()
// 	predFile      = app.Flag("predfile", "Predictor JSON file.").Required().String()
// 	takeProfit    = app.Flag("takeprofit", "Take profit level in price points. Default is 1").Float64()
// 	stopLoss      = app.Flag("stoploss", "Stop loss level in price points. Default is 0.5").Float64()
// 	savePositions = app.Flag("savepositions", "Save positions to a JSONL file. Default is false.").Bool()
// )

// const (
// 	DefaultTakeProfit = 1.0
// 	DefaultStopLoss   = 0.5
// )

// func main() {
// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

// 	zerolog.SetGlobalLevel(zerolog.InfoLevel)

// 	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

// 	if *debug {
// 		zerolog.SetGlobalLevel(zerolog.DebugLevel)
// 	}

// 	if *takeProfit == 0.0 {
// 		*takeProfit = DefaultTakeProfit
// 	}

// 	if *stopLoss == 0.0 {
// 		*stopLoss = DefaultStopLoss
// 	}

// 	s, err := collection.NewDirStore(*dataDir)
// 	if err != nil {
// 		log.Fatal().Err(err).Str("dataDir", *dataDir).Msg("failed to make dir store")
// 	}

// 	c, err := collection.Load(s)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to load collection")
// 	}

// 	pred, err := predictors.LoadPredictorFromFile(*predFile)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to load predictor")
// 	}

// 	log.Info().Str("params", pred.Params().String()).Msg("loaded predictor")

// 	positions := models.Positions{}

// 	evalPredictions := func(dir models.Direction, bars models.Bars) {
// 		if pos := evalPrediction(dir, bars); pos != nil {
// 			positions = append(positions, pos)
// 		}
// 	}

// 	dc := backtesting.NewSequentialDates(c.GetTimeSpan().DateRangeIn(time.Local))
// 	dt := backtesting.NewDayTrader(c, pred, dc, evalPredictions)

// 	fmt.Println("Running backtests")
// 	for dt.AnyLeft() {
// 		fmt.Printf(".")

// 		err := dt.RunTest()
// 		if err != nil {
// 			log.Fatal().Err(err).Msg("backtest failed")
// 		}

// 		dt.Advance()
// 	}

// 	fmt.Printf("\n\n# Trades: %d\n", len(positions))

// 	analysis := positions.Analyze()

// 	fmt.Printf("Winning trades: %f%%\n", analysis.Winning*100.0)
// 	fmt.Printf("Total P/L: %f\n", analysis.TotalPL)

// 	if *savePositions {
// 		storePositions(positions)
// 	}
// }

// func storePositions(positions models.Positions) {
// 	fpath := "./positions.jsonl"

// 	positions.Localize()

// 	err := positions.StoreToFile(fpath)
// 	if err != nil {
// 		log.Error().Err(err).Msg("failed to save positions to file")

// 		return
// 	}

// 	log.Info().Str("path", fpath).Msg("saved positions")
// }

// func evalPrediction(dir models.Direction, bars models.Bars) *models.Position {
// 	if len(bars) < 2 {
// 		log.Fatal().Int("len", len(bars)).Msg("not enough bars to evaluate prediction")
// 	}

// 	var pos *models.Position

// 	switch dir {
// 	case models.DirDown:
// 		pos = models.NewShortPosition(bars[0].Timestamp, bars[0].Close)
// 	case models.DirUp:
// 		pos = models.NewLongPosition(bars[0].Timestamp, bars[0].Close)
// 	case models.DirNone:
// 		return nil
// 	default:
// 		log.Fatal().Int("dir", int(dir)).Msg("unknown direction")
// 	}

// 	for i := 1; pos.IsOpen() && (i < len(bars)); i++ {
// 		bar := bars[i]
// 		pl, _ := pos.OpenProfitLoss(bar.Close)

// 		if pl < -*stopLoss {
// 			pos.Close(bar.Timestamp, bar.Close, "stop loss")
// 		} else if pl > *takeProfit {
// 			pos.Close(bar.Timestamp, bar.Close, "take profit")
// 		}
// 	}
// 	if pos.IsOpen() {
// 		lastBar := bars.Last()

// 		pos.Close(lastBar.Timestamp, lastBar.Close, "change dir/end test")
// 	}

// 	return pos
// }
