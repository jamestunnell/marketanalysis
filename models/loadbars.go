package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func LoadBars(fpath string) ([]*Bar, error) {
	f, err := os.Open(fpath)
	if err != nil {
		err = fmt.Errorf("failed to open file %s: %w", fpath, err)

		return []*Bar{}, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	bars := make([]*Bar, 0)
	line := 1

	for scanner.Scan() {
		var bar Bar

		if err := json.Unmarshal([]byte(scanner.Text()), &bar); err != nil {
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

	return bars, err
}
