package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
)

type config struct {
	TimeZone int
}

const (
	localConfigFilename = "config.json"
)

func main() {
	cfg, err := loadOrCreateLocalConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create/load config")
	}

	a := app.New()
	w := a.NewWindow("Market Data Manager")

	tz := widget.NewEntry()
	tz.SetText(strconv.Itoa(cfg.TimeZone))

	setupContent := container.NewVBox(
		container.NewHBox(widget.NewLabel("Time Zone"), tz),
	)
	setupTab := container.NewTabItem("Setup", setupContent)

	appTabs := container.NewAppTabs(setupTab)

	w.SetContent(appTabs)
	w.ShowAndRun()
}

func localConfigRoot() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get user home dir")
	}

	return filepath.Join(dirname, ".local/share")
}

func loadOrCreateLocalConfig() (*config, error) {
	cfgRoot := localConfigRoot()

	_, err := os.Stat(cfgRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to stat local config root '%s': %w", cfgRoot, err)
	}

	localConfigPath := filepath.Join(cfgRoot, localConfigFilename)

	_, err = os.Stat(localConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			d, err := json.Marshal(makeDefaultConfig())
			if err != nil {
				return nil, fmt.Errorf("failed to marshal default config: %w", err)
			}

			if err = os.WriteFile(localConfigPath, d, 0644); err != nil {
				return nil, fmt.Errorf("failed to write default config to '%s': %w", localConfigPath, err)
			}
		} else {
			return nil, fmt.Errorf("failed to stat local config file '%s': %w", localConfigPath, err)
		}
	}

	d, err := os.ReadFile(localConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", localConfigPath, err)
	}

	var cfg config

	if err = json.Unmarshal(d, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &cfg, nil
}

func makeDefaultConfig() *config {
	return &config{
		TimeZone: 0,
	}
}
