package collectbars

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rs/zerolog/log"
)

type Params struct {
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Outpath string    `json:"outpath"`
	Symbol  string    `json:"symbol"`
}

type Command struct {
	*Params
}

var (
	errEndBeforeStart = errors.New("end time is before start")
	errNoSymbol       = errors.New("no symbol given")
	errOutdirNotExist = errors.New("outdir does not exist")
)

func New(params *Params) (*Command, error) {
	if params.Symbol == "" {
		return nil, errNoSymbol
	}

	if params.End.Before(params.Start) {
		return nil, errEndBeforeStart
	}

	_, err := os.Stat(filepath.Dir(params.Outpath))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errOutdirNotExist
		} else {
			return nil, fmt.Errorf("failed to stat outpath: %w", err)
		}
	}

	cmd := &Command{Params: params}

	return cmd, nil
}

func (params *Params) Format() string {
	d, err := json.Marshal(params)
	if err != nil {
		log.Warn().Err(err).Msg("failed to marshal collectbars.Params")

		return ""
	}

	return string(d)
}

func (cmd *Command) Run() error {
	fmt.Printf("params: %s\n", cmd.Format())

	fmt.Printf("start time: %s\n", cmd.Start.Format(time.RFC3339))
	fmt.Printf("end time: %s\n", cmd.End.Format(time.RFC3339))

	fmt.Printf("creating JSON Lines data file %s\n", cmd.Outpath)
	outFile, err := os.Create(cmd.Outpath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", cmd.Outpath, err)
	}

	fmt.Printf("collecting bars for %s\n", cmd.Symbol)

	alpacaBars, err := marketdata.GetBars(cmd.Symbol, marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneMin,
		Start:     cmd.Start,
		End:       cmd.End,
		AsOf:      "-",
	})
	if err != nil {
		return fmt.Errorf("failed to get bars: %w\n", err)
	}

	fmt.Printf("collected %d bars\n", len(alpacaBars))

	fmt.Println("writing bar data to file")

	totalBytes := 0

	for _, alpacaBar := range alpacaBars {
		bar := models.NewFromAlpacaBar(alpacaBar)

		d, err := json.Marshal(bar)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON for bar\n: - bar value: %v\n - error: %w", bar, err)
		}

		n, err := outFile.Write(d)
		if err != nil {
			return fmt.Errorf("failed to write bar data to file: %w", err)
		}

		totalBytes += n

		n, err = outFile.WriteString("\n")
		if err != nil {
			return fmt.Errorf("failed to write newline delimiter to file: %w", err)
		}

		totalBytes += n
	}

	if err = outFile.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	fmt.Printf("wrote %d bytes\n", totalBytes)

	return nil
}
