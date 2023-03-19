package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rs/zerolog/log"
)

var (
	app = kingpin.New("plot", "Plot market data.")

	bars = app.Command("bars", "Plot bar data.")

	barsIn      = bars.Flag("in", "Input file").Required().String()
	barsHTMLOut = bars.Flag("htmlout", "HTML output file").Required().String()
	barsSeries  = bars.Flag("series", "Data series name").String()
)

func klineChart(title, seriesName string, bars []*models.Bar) *charts.Kline {
	if seriesName == "" {
		seriesName = "data series"
	}

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
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case bars.FullCommand():
		bars, err := models.LoadBars(*barsIn)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load bars")
		}

		title := filepath.Base(*barsIn)
		chart := klineChart(title, *barsSeries, bars)

		f, err := os.Create(*barsHTMLOut)
		if err != nil {
			log.Fatal().Err(err).Str("htmlout", *barsHTMLOut).Msg("failed to create HTML output file")
		}

		chart.Render(f)
	}

	// 	http.HandleFunc("/", httpserver)
	// 	http.ListenAndServe(":8081", nil)
}
