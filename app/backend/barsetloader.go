package backend

import (
	"context"
	"fmt"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/loading"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketdata"
	"github.com/jamestunnell/marketdata/alpaca"
)

type BarSetLoader struct {
	Symbol string
	Store  Store[*models.BarSet]
}

var locNY = loading.GetLocationNY()

func NewBarSetLoader(
	db *mongo.Database,
	symbol string,
) *BarSetLoader {
	info := &ResourceInfo{
		KeyName:    "date",
		Name:       "barset",
		NamePlural: "barsets",
	}
	col := db.Collection(symbol)
	store := NewMongoStore[*models.BarSet](info, col)

	return &BarSetLoader{
		Symbol: symbol,
		Store:  store,
	}
}

func (l *BarSetLoader) Load(ctx context.Context, d date.Date) (marketdata.Bars, error) {
	log.Trace().Stringer("date", d).Msg("loading day bars")

	dayBars, appErr := l.Store.Get(ctx, d.String())
	if appErr == nil {
		log.Trace().
			Int("count", len(dayBars.Bars)).
			Stringer("date", d).
			Msg("found bars in store")

		return dayBars.Bars, nil
	}

	// return marketdata.Bars{}, fmt.Errorf("failed to get day bars: %w", appErr)

	ts := timespan.NewTimeSpan(d.In(locNY), d.Add(1).In(locNY))
	bc := alpaca.NewFreeBarCollector()

	bars, err := bc.Collect(l.Symbol, ts)
	if err != nil {
		return nil, fmt.Errorf("failed to get aplaca bars: %w", err)
	}

	for _, bar := range bars {
		bar.Timestamp = bar.Timestamp.In(locNY)
	}

	log.Trace().
		Int("count", len(bars)).
		Stringer("date", d).
		Msg("loaded bars from alpaca")

	dayBars = &models.BarSet{
		Bars: bars,
		Date: d.String(),
	}

	if d.Equal(date.TodayIn(locNY)) {
		log.Trace().Msg("not storing bars from today")

		return dayBars.Bars, nil
	}

	log.Debug().
		Int("count", len(bars)).
		Stringer("date", d).
		Msg("storing day bars")

	appErr = l.Store.Create(ctx, dayBars)
	if appErr != nil {
		log.Warn().
			Err(appErr).
			Stringer("date", d).
			Msg("failed to store day bars")
	} else {
		log.Debug().
			Int("count", len(bars)).
			Stringer("date", d).
			Msg("stored bars")
	}

	return dayBars.Bars, nil
}
