package processing

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/provision"
)

type EvalStepFunc func(bar *models.Bar, sourceOut, procsOut float64) error

func Evaluate(chain *Chain, bars provision.BarSequence, step EvalStepFunc) error {
	if err := chain.Initialize(); err != nil {
		return fmt.Errorf("failed to init chain: %w", err)
	}

	eachBar := func(bar *models.Bar) error {
		chain.Update(bar)

		if chain.SourceWarm() && chain.ProcsWarm() {
			err := step(bar, chain.SourceOutput(), chain.ProcsOutput())
			if err != nil {
				return fmt.Errorf("eval step failed: %w", err)
			}
		}

		return nil
	}

	if err := bars.EachBar(eachBar); err != nil {
		return fmt.Errorf("each bar failed: %w", err)
	}

	return nil
}
