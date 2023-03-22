package datamanager

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	Resolution1Min = "1m"
)

type CollectionInfo struct {
	Symbol     string `json:"symbol"`
	Resolution string `json:"resolution"`
}

func NewCollectionInfo(sym, res string) *CollectionInfo {
	return &CollectionInfo{
		Symbol:     sym,
		Resolution: res,
	}
}

func (info *CollectionInfo) Store(fpath string) error {
	d, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err = os.WriteFile(fpath, d, 0644); err != nil {
		return fmt.Errorf("failed to write file '%s': %w", fpath, err)
	}

	return nil
}

func (info *CollectionInfo) Load(fpath string) error {
	d, err := os.ReadFile(fpath)
	if err != nil {
		return fmt.Errorf("failed to read file '%s': %w", fpath, err)
	}

	if err = json.Unmarshal(d, info); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}
