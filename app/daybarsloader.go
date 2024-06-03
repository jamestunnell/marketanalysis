package app

import (
	"context"
	"fmt"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/models"
)

type DayBarsLoader struct {
	Symbol   string
	Location *time.Location
	Store    Store[*models.DayBars]
}

func NewDayBarsLoader(
	db *mongo.Database,
	symbol string,
	loc *time.Location,
) models.DayBarsLoader {
	info := &ResourceInfo{
		KeyName:    "date",
		Name:       "barset",
		NamePlural: "barsets",
	}
	col := db.Collection(symbol)
	store := NewMongoStore[*models.DayBars](info, col)

	return &DayBarsLoader{
		Symbol:   symbol,
		Location: loc,
		Store:    store,
	}
}

func (l *DayBarsLoader) Load(ctx context.Context, d date.Date) (*models.DayBars, error) {
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

		return dayBars, nil
	}

	ts := timespan.NewTimeSpan(d.In(l.Location), d.Add(1).In(l.Location))

	bs, err := bars.GetAlpacaBarsOneMin(l.Symbol, ts, l.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to get aplaca bars: %w", err)
	}

	log.Debug().
		Int("count", len(bs)).
		Stringer("date", d).
		Msg("loaded bars from alpaca")

	dayBars = &models.DayBars{
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

	return dayBars, nil
}
