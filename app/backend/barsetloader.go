package backend

import (
	"context"
	"fmt"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/alpaca"
	"github.com/jamestunnell/marketanalysis/models"
)

type BarSetLoader struct {
	Symbol string
	Store  Store[*models.BarSet]
}

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

func (l *BarSetLoader) Load(ctx context.Context, d date.Date) (models.Bars, error) {
	log.Debug().Stringer("date", d).Msg("loading day bars")

	dayBars, appErr := l.Store.Get(ctx, d.String())
	if appErr == nil {
		log.Debug().
			Int("count", len(dayBars.Bars)).
			Stringer("date", d).
			Msg("found bars in store")

		return dayBars.Bars, nil
	}

	ts := timespan.NewTimeSpan(d.In(locNewYork), d.Add(1).In(locNewYork))

	bs, err := alpaca.GetBarsOneMin(l.Symbol, ts, locNewYork)
	if err != nil {
		return nil, fmt.Errorf("failed to get aplaca bars: %w", err)
	}

	log.Debug().
		Int("count", len(bs)).
		Stringer("date", d).
		Msg("loaded bars from alpaca")

	dayBars = &models.BarSet{
		Bars: bs,
		Date: d.String(),
	}

	if d.Equal(date.TodayIn(locNewYork)) {
		log.Debug().Msg("not storing bars from today")

		return dayBars.Bars, nil
	}

	log.Debug().
		Int("count", len(bs)).
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
			Int("count", len(bs)).
			Stringer("date", d).
			Msg("stored bars")
	}

	return dayBars.Bars, nil
}
