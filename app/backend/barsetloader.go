package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/alpaca"
	"github.com/jamestunnell/marketanalysis/models"
)

type BarSetLoader struct {
	Symbol   string
	Location *time.Location
	Store    Store[*models.BarSet]
}

func NewBarSetLoader(
	db *mongo.Database,
	symbol string,
	loc *time.Location,
) *BarSetLoader {
	info := &ResourceInfo{
		KeyName:    "date",
		Name:       "barset",
		NamePlural: "barsets",
	}
	col := db.Collection(symbol)
	store := NewMongoStore[*models.BarSet](info, col)

	return &BarSetLoader{
		Symbol:   symbol,
		Location: loc,
		Store:    store,
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

		for _, bar := range dayBars.Bars {
			bar.Timestamp = bar.Timestamp.In(l.Location)
		}

		return dayBars.Bars, nil
	}

	ts := timespan.NewTimeSpan(d.In(l.Location), d.Add(1).In(l.Location))

	bs, err := alpaca.GetBarsOneMin(l.Symbol, ts, l.Location)
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
