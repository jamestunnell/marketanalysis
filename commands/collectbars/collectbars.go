package collectbars

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/rs/zerolog/log"
)

type Params struct {
	Start   string `json:"start"`
	Dur     string `json:"dur"`
	Outdir  string `json:"outdir"`
	Outfile string `json:"outfile"`
	Symbol  string `json:"symbol"`
}

type Command struct {
	Symbol     string
	Start, End time.Time
	OutPath    string
}

var errNoSymbol = errors.New("no symbol given")

func New(params *Params) (*Command, error) {
	params.SetDefaults()

	if params.Symbol == "" {
		return nil, errNoSymbol
	}

	fmt.Printf("params: %s\n", params.Format())

	t, err := time.Parse(time.RFC3339, params.Start)
	if err != nil {
		err = fmt.Errorf("failed to parse start datetime '%s' using RFC3339: %w", params.Start, err)
	}

	dur, err := time.ParseDuration(params.Dur)
	if err != nil {
		err = fmt.Errorf("failed to parse duration '%s': %w", params.Dur, err)
	}

	outPath := filepath.Join(params.Outdir, params.Outfile)

	cmd := &Command{
		Symbol:  params.Symbol,
		Start:   t,
		End:     t.Add(dur),
		OutPath: outPath,
	}

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

func (params *Params) SetDefaults() {
	if params.Dur == "" {
		params.Dur = "30m"
	}

	if params.Start == "" {
		t := time.Now().Add(-time.Hour)

		params.Start = t.Format(time.RFC3339)
	}

	if params.Outdir == "" {
		params.Outdir = "."
	}

	if params.Outfile == "" {
		params.Outfile = "bars.jsonl"
	}
}

func (cmd *Command) Run() error {
	fmt.Printf("start time: %s\n", cmd.Start.Format(time.RFC3339))
	fmt.Printf("end time: %s\n", cmd.End.Format(time.RFC3339))

	fmt.Printf("creating JSON Lines data file %s\n", cmd.OutPath)
	outFile, err := os.Create(cmd.OutPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", cmd.OutPath, err)
	}

	fmt.Printf("collecting bars for %s\n", cmd.Symbol)

	bars, err := marketdata.GetBars(cmd.Symbol, marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneMin,
		Start:     cmd.Start,
		End:       cmd.End,
		AsOf:      "-",
	})
	if err != nil {
		return fmt.Errorf("failed to get bars: %w\n", err)
	}

	fmt.Printf("collected %d bars\n", len(bars))

	fmt.Println("writing bar data to file")

	totalBytes := 0

	for _, bar := range bars {
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
