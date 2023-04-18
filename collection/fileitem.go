package collection

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileItem struct {
	name, fpath string
}

func NewFileItem(fpath string) (Item, error) {
	info, err := os.Stat(fpath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat root dir '%s': %w", fpath, err)
	}

	if info.IsDir() {
		return nil, &ErrNotFile{Path: fpath}
	}

	item := &FileItem{
		name:  filepath.Base(fpath),
		fpath: fpath,
	}

	return item, nil
}

func (item *FileItem) Name() string {
	return item.name
}

func (item *FileItem) Load() ([]byte, error) {
	d, err := os.ReadFile(item.fpath)
	if err != nil {
		err = fmt.Errorf("failed to read item '%s': %w", item.name, err)
		return []byte{}, err
	}

	return d, nil
}

func (item *FileItem) Store(data []byte) error {
	if err := os.WriteFile(item.fpath, data, 0644); err != nil {
		return fmt.Errorf("failed to read item '%s': %w", item.name, err)
	}

	return nil
}
