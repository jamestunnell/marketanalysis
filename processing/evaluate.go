package processing

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/provision"
)

type EvalStepFunc func(bar *models.Bar, sourceOut, procsOut float64)

func Evaluate(chain *Chain, bars provision.BarSequence, step EvalStepFunc) error {
	if err := chain.Initialize(); err != nil {
		return fmt.Errorf("failed to init chain: %w", err)
	}

	if err := bars.Initialize(); err != nil {
		return fmt.Errorf("failed to init bar sequence: %w", err)
	}

	bars.EachBar(func(bar *models.Bar) {
		chain.Update(bar)

		if chain.SourceWarm() && chain.ProcsWarm() {
			step(bar, chain.SourceOutput(), chain.ProcsOutput())
		}
	})

	return nil
}
