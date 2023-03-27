package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func StoreBars(bars []*Bar, fpath string) error {
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fpath, err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	for i, bar := range bars {
		d, err := json.Marshal(bar)
		if err != nil {
			return fmt.Errorf("failed to marshal bar %d: %w", i+1, err)
		}

		_, err = w.Write(d)
		if err != nil {
			return fmt.Errorf("failed to write bar %d: %w", i+1, err)
		}
	}

	return err
}
