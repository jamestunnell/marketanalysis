package collectbars

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rs/zerolog/log"
)

type Params struct {
	Start         time.Time `json:"start"`
	End           time.Time `json:"end"`
	CollectionDir string    `json:"collectionDir"`
	Symbol        string    `json:"symbol"`
}

type Command struct {
	*Params
	Store collection.Store
}

var (
	errEndBeforeStart = errors.New("end time is before start")
	errNoSymbol       = errors.New("no symbol given")
)

func New(params *Params) (*Command, error) {
	const fmtNewStoreFailed = "failed to make store for collection dir '%s': %w"

	if params.Symbol == "" {
		return nil, errNoSymbol
	}

	if params.End.Before(params.Start) {
		return nil, errEndBeforeStart
	}

	store, err := collection.NewDirStore(params.CollectionDir)
	if err != nil {
		err = fmt.Errorf(fmtNewStoreFailed, params.CollectionDir, err)

		return nil, err
	}

	cmd := &Command{
		Params: params,
		Store:  store,
	}

	return cmd, nil
}

func (params *Params) FormatParams() string {
	d, err := json.Marshal(params)
	if err != nil {
		log.Warn().Err(err).Msg("failed to marshal collectbars.Params")

		return ""
	}

	return string(d)
}

func (cmd *Command) Run() error {
	fmt.Printf("params: %s\n", cmd.FormatParams())

	fmt.Printf("start time: %s\n", cmd.Start.Format(time.RFC3339))
	fmt.Printf("end time: %s\n", cmd.End.Format(time.RFC3339))

	alpacaBars, err := marketdata.GetBars(cmd.Symbol, marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneMin,
		Start:     cmd.Start,
		End:       cmd.End,
		AsOf:      "-",
	})
	if err != nil {
		return fmt.Errorf("failed to get bars: %w", err)
	}

	fmt.Printf("collected %d bars\n", len(alpacaBars))

	bars := make([]*models.Bar, len(alpacaBars))
	for i, alpacaBar := range alpacaBars {
		bars[i] = models.NewFromAlpacaBar(alpacaBar)
	}

	exists, err := collection.Exists(cmd.Store)
	if err != nil {
		return fmt.Errorf("failed to check if collection exists: %w", err)
	}

	var c collection.Collection

	if exists {
		c, err = collection.Load(cmd.Store)
		if err != nil {
			return fmt.Errorf("failed to load collection: %w", err)
		}

		sym := c.Info().Symbol
		if sym != cmd.Symbol {
			err = fmt.Errorf(
				"collection symbol '%s' does not match given '%s'", sym, cmd.Symbol)

			return err
		}

		added := c.AddBars(bars)

		fmt.Printf("added %d bars to existing collection", added)
	} else {
		info := collection.NewInfo(cmd.Symbol, collection.Resolution1Min)

		c, err = collection.New(info, bars)
		if err != nil {
			return fmt.Errorf("failed to create new collection: %w", err)
		}

		fmt.Println("created new collection")
	}

	if err = c.Store(cmd.Store); err != nil {
		return fmt.Errorf("failed to store collection: %w", err)
	}

	return nil
}
