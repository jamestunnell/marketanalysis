package backtesting

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
)

func Backtest(s models.Strategy, bars models.Bars) error {
	if len(bars) <= s.WarmupPeriod() {
		return commonerrs.NewErrLessThanMin("warmup+backtest bars", s.WarmupPeriod()+1, len(bars))
	}

	wuBars := bars[:s.WarmupPeriod()]
	remBars := bars[s.WarmupPeriod():]

	if err := s.Initialize(wuBars); err != nil {
		return fmt.Errorf("failed to warm up strategy")
	}

	for _, bar := range remBars {
		s.Update(bar)
	}

	s.Close(bars[len(bars)-1], "end-of-backtest")

	return nil
}
