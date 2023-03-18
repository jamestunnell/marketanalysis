package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/commands/collectbars"
)

type Command interface {
	Run() error
}

var (
	app = kingpin.New("collect", "Collect market data.")

	bars = app.Command("bars", "Collect historical 1-minute bar data.")

	barsStart   = bars.Flag("start", "Start date-time formated according to RFC3339. Default is 1 hour ago.").String()
	barsDur     = bars.Flag("dur", "Collection period duration in hours and/or minutes (e.g.: 10h, 1h10m, 50m). Default is 30m.").String()
	barsOutdir  = bars.Flag("outdir", "Output directory for bar data file. Default is current working dir.").String()
	barsOutfile = bars.Flag("outfile", "Output filename for bar data file. Default is bars.jsonl.").String()
	barsSymbol  = bars.Arg("symbol", "The stock symbol.").String()
)

func main() {
	var cmd Command

	var err error

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case bars.FullCommand():
		params := &collectbars.Params{
			Start:   *barsStart,
			Dur:     *barsDur,
			Outdir:  *barsOutdir,
			Outfile: *barsOutfile,
			Symbol:  *barsSymbol,
		}
		cmd, err = collectbars.New(params)
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
