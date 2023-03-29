package bar

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func StoreBarsToFile(bars []*Bar, fpath string) error {
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fpath, err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	return StoreBars(bars, w)
}

func StoreBars(bars []*Bar, w io.Writer) error {
	for i, bar := range bars {
		d, err := json.Marshal(bar)
		if err != nil {
			return fmt.Errorf("failed to marshal bar %d: %w", i+1, err)
		}

		d = append(d, byte('\n'))

		_, err = w.Write(d)
		if err != nil {
			return fmt.Errorf("failed to write bar %d: %w", i+1, err)
		}
	}

	return nil
}
