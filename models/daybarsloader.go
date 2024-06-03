package models

import (
	"context"

	"github.com/rickb777/date"
)

type DayBarsLoader interface {
	Load(ctx context.Context, d date.Date) (*DayBars, error)
}
