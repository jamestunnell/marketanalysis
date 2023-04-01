package main

import (
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models/bar"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
)

var (
	app = kingpin.New("plot", "Plot market data.")

	dataDir = app.Flag("datadir", "Data collection root dir").Required().String()
	htmlOut = app.Flag("htmlout", "HTML output file").Required().String()
	start   = app.Flag("start", "start datetime formatted in RFC3339").String()
	end     = app.Flag("end", "end datetime formatted in RFC3339").String()

	barsCmd = app.Command("bars", "Plot bar data.")
	atrCmd  = app.Command("atr", "Plot ATR.")

	atrLength = atrCmd.Flag("n", "length, must be positive").Required().Int()
)

func lineChart(title, seriesName string, times []time.Time, data []float64) *charts.Line {
	line := charts.NewLine()

	if len(times) != len(data) {
		log.Fatal().Msg("len mismatch")
	}

	x := make([]string, 0)
	y := make([]opts.LineData, 0)
	for i := 0; i < len(times); i++ {
		x = append(x, times[i].Format(time.RFC3339))
		y = append(y, opts.LineData{Value: data[i]})
	}

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title}),
	)

	line.SetXAxis(x).
		AddSeries(seriesName, y)

	return line
}

func klineChart(title, seriesName string, bars []*bar.Bar) *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(bars); i++ {
		x = append(x, bars[i].Timestamp.Format(time.RFC3339))
		y = append(y, opts.KlineData{Value: bars[i].OpenCloseLowHigh()})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries(seriesName, y)

	return kline
}

// func httpserver(w http.ResponseWriter, _ *http.Request) {
// 	page := components.NewPage()
// 	page.AddCharts(
// 		klineBase(),
// 	)
// 	page.Render(w)
// }

func main() {
	cmdName := kingpin.MustParse(app.Parse(os.Args[1:]))

	s, err := collection.NewDirStore(*dataDir)
	if err != nil {
		log.Fatal().Err(err).Str("dataDir", *dataDir).Msg("failed to make dir store")
	}

	c, err := collection.Load(s)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load collection")
	}

	ts := c.Timespan()
	tStart := ts.Start()
	tEnd := ts.End()

	if *start != "" {
		tStart, err = time.Parse(time.RFC3339, *start)
		if err != nil {
			log.Fatal().Err(err).Str("start", *start).Msg("failed to parse start time")
		}
	}

	if *end != "" {
		tEnd, err = time.Parse(time.RFC3339, *end)
		if err != nil {
			log.Fatal().Err(err).Str("end", *end).Msg("failed to parse end time")
		}
	}

	bars := c.GetBars(timespan.NewTimeSpan(tStart, tEnd))

	switch cmdName {
	case barsCmd.FullCommand():
		chart := klineChart(c.Info().Symbol, "bar data", bars)

		f, err := os.Create(*htmlOut)
		if err != nil {
			log.Fatal().Err(err).Str("htmlout", *htmlOut).Msg("failed to create HTML output file")
		}

		chart.Render(f)
	case atrCmd.FullCommand():
		atr := indicators.NewATR(*atrLength)
		warmupBars := bars[:atr.WarmupPeriod()]
		remainingBars := bars[atr.WarmupPeriod():]

		err = atr.Initialize(warmupBars)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to init ATR indicator")
		}

		atrVals := make([]float64, len(remainingBars))
		times := make([]time.Time, len(remainingBars))
		for i, bar := range remainingBars {
			times[i] = bar.Timestamp
			atrVals[i] = atr.Update(bar)
		}

		chart := lineChart(c.Info().Symbol, "bar data", times, atrVals)

		f, err := os.Create(*htmlOut)
		if err != nil {
			log.Fatal().Err(err).Str("htmlout", *htmlOut).Msg("failed to create HTML output file")
		}

		chart.Render(f)
	}

	// 	http.HandleFunc("/", httpserver)
	// 	http.ListenAndServe(":8081", nil)
}
