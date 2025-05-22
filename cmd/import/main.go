package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/rickb777/date"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"

	"github.com/jamestunnell/marketanalysis/app/backend/appvars"
	"github.com/jamestunnell/marketanalysis/app/backend/database"
	"github.com/jamestunnell/marketanalysis/app/backend/stores"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketdata"
)

func main() {
	inFile := kingpin.Flag("in", "Input tar.gz file.").Required().String()
	sym := kingpin.Flag("sym", "The stock symbol.").Required().String()

	_ = kingpin.Parse()

	vars, err := appvars.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load app vars")
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if vars.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	client, err := database.ConnectMongo(vars)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to mongo DB")
	}

	defer database.DisconnectMongo(client)

	db := client.Database(database.MongoDBName)
	store := stores.NewBarSets(*sym, db)
	eachFile := func(r io.Reader) {
		bars, err := marketdata.LoadBars(r)
		if err != nil {
			log.Warn().Err(err).Msg("failed to load bars: %w")

			return
		}

		numImported := 0

		for _, bs := range makeBarSets(bars) {
			err = store.Upsert(context.Background(), bs)
			if err != nil {
				log.Warn().Err(err).Msg("failed to upsert bar set: %w")

				continue
			}

			numImported += len(bs.Bars)
		}

		log.Info().Str("sym", *sym).Int("bars", numImported).Msg("imported bar data")
	}

	if err := upackTarGz(*inFile, eachFile); err != nil {
		log.Fatal().Err(err).Msg("failed to unpack tar.gz file")

		return
	}

	log.Info().Msg("imported complete")
}

func upackTarGz(fpath string, eachFile func(io.Reader)) error {
	f, err := os.Open(fpath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", fpath, err)
	}

	zr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to make gzip reader: %w", err)
	}

	defer func() {
		if closeErr := zr.Close(); closeErr != nil {
			log.Warn().Err(err).Msg("failed to close gzip reader")
		}
	}()

	var tarData []byte

	if tarData, err = io.ReadAll(zr); err != nil {
		return fmt.Errorf("failed to read gzip data: %w", err)
	}

	tr := tar.NewReader(bytes.NewReader(tarData))

	for {
		_, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}

		if err != nil {
			return fmt.Errorf("failed to read next archive file: %w", err)
		}

		var buf bytes.Buffer

		if _, err := io.Copy(&buf, tr); err != nil {
			return fmt.Errorf("failed to get read archive file: %w", err)
		}

		eachFile(&buf)
	}

	return nil
}

func makeBarSets(bars marketdata.Bars) []*models.BarSet {
	barSets := map[string]*models.BarSet{}

	for _, b := range bars {
		dateStr := b.Date().Format(date.RFC3339)

		bs, found := barSets[dateStr]
		if found {
			bs.Bars = append(bs.Bars, b)
		} else {
			bs = &models.BarSet{
				Date: dateStr, Bars: marketdata.Bars{b},
			}

			barSets[dateStr] = bs
		}
	}

	return maps.Values(barSets)
}
