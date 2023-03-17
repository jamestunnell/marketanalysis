package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

type Command interface {
	Run() error
}

type CollectDay struct {
	Symbol string
	Date   time.Time
}

var (
	app = kingpin.New("collectbars", "A command-line bar data collection.")
	//debug    = app.Flag("debug", "Enable debug mode.").Bool()
	//serverIP = app.Flag("server", "Server address.").Default("127.0.0.1").IP()

	day       = app.Command("day", "Collect a whole days worth of bars.")
	daySymbol = day.Arg("symbol", "The stock symbol.").Required().String()
	dayDate   = day.Arg("date", "Date formated YYYY-MM-DD.").Required().String()
)

func main() {
	var cmd Command

	var err error

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case day.FullCommand():
		cmd, err = NewCollectDay(*daySymbol, *dayDate)
	}

	if err != nil {
		panic(fmt.Errorf("failed to make command: %w", err))
	}

	if cmd == nil {
		panic("unknown command")
	}

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func NewCollectDay(symbol, dateStr string) (*CollectDay, error) {
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date '%s'", dateStr)
	}

	cmd := &CollectDay{
		Symbol: symbol,
		Date:   date,
	}

	return cmd, nil
}

func (cmd *CollectDay) Run() error {
	year, month, day := cmd.Date.Date()

	fmt.Printf("collecting bars for %s on %04d-%02d-%02d\n", cmd.Symbol, year, month, day)

	bars, err := marketdata.GetBars(cmd.Symbol, marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneMin,
		Start:     time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		End:       time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC),
		AsOf:      "-",
	})
	if err != nil {
		return fmt.Errorf("failed to get bars: %w\n", err)
	}

	fmt.Printf("collected %d bars", len(bars))

	for _, bar := range bars {
		fmt.Printf("%+v\n", bar)
	}

	return nil
}
