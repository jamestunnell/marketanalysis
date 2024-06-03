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
	dayBars, appErr := l.Store.Get(ctx, d.String())
	if appErr == nil {
		return dayBars, nil
	}

	ts := timespan.NewTimeSpan(d.In(l.Location), d.Add(1).In(l.Location))

	bars, err := bars.GetAlpacaBarsOneMin(l.Symbol, ts, l.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to get aplaca bars: %w", err)
	}

	dayBars = &models.DayBars{
		Bars: bars,
		Date: d,
	}

	appErr = l.Store.Create(ctx, dayBars)
	if appErr != nil {
		log.Warn().Err(err).Stringer("date", d).Msg("failed to store day bars")
	}

	return dayBars, nil
}
