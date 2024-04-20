package collect

import (
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rs/zerolog/log"
)

func GetAlpacaBars(start, end time.Time, sym string, loc *time.Location) (models.Bars, error) {
	// the most current end time alpaca allows for free
	latestEndAllowed := time.Now().Add(-15 * time.Minute)
	if end.After(latestEndAllowed) {
		end = latestEndAllowed
	}

	alpacaBars, err := marketdata.GetBars(sym, marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneMin,
		Start:     start,
		End:       end,
		AsOf:      "-",
	})
	if err != nil {
		return models.Bars{}, fmt.Errorf("failed to get bars from alpaca: %w", err)
	}

	log.Info().
		Time("start", start).
		Time("end", end).
		Int("count", len(alpacaBars)).
		Msg("collected bars from alpaca")

	bars := make([]*models.Bar, len(alpacaBars))
	for i, alpacaBar := range alpacaBars {
		bar := models.NewBarFromAlpaca(alpacaBar)

		bar.Timestamp = bar.Timestamp.In(loc)

		bars[i] = bar
	}

	return bars, nil
}