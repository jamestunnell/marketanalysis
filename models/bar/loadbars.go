package bar

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func LoadBarsFromFile(fpath string) ([]*Bar, error) {
	f, err := os.Open(fpath)
	if err != nil {
		err = fmt.Errorf("failed to open file %s: %w", fpath, err)

		return []*Bar{}, err
	}

	defer f.Close()

	return LoadBars(f)
}

func LoadBars(r io.Reader) ([]*Bar, error) {
	scanner := bufio.NewScanner(r)
	bars := make([]*Bar, 0)
	line := 1

	for scanner.Scan() {
		var bar Bar

		lineStr := scanner.Text()

		// ignore empty lines
		if lineStr == "" {
			continue
		}

		if err := json.Unmarshal([]byte(lineStr), &bar); err != nil {
			err = fmt.Errorf("failed to unmarshal line %d :%w", line, err)

			return []*Bar{}, err
		}

		bars = append(bars, &bar)

		line++
	}

	if err := scanner.Err(); err != nil {
		err = fmt.Errorf("scanner failed: %w", err)

		return []*Bar{}, err
	}

	return bars, nil
}
