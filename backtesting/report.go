package backtesting

import (
	"time"

	"github.com/jamestunnell/marketanalysis/models"
)

type Report struct {
	Start     time.Time
	Positions models.Positions
}
