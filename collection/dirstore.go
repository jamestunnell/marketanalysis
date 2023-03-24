package collection

import (
	"fmt"
	"os"
	"path/filepath"
)

type DirStore struct {
	root string
}

func NewDirStore(root string) (Store, error) {
	rootInfo, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("failed to stat root dir '%s': %w", root, err)
	}

	if !rootInfo.IsDir() {
		return nil, &ErrNotDir{DirPath: root}
	}

	store := &DirStore{root: root}

	return store, nil
}

func (store *DirStore) ItemNames() ([]string, error) {
	entries, err := os.ReadDir(store.root)
	if err != nil {
		err = fmt.Errorf("failed to read root dir '%s': %w", store.root, err)

		return []string{}, err
	}

	names := []string{}

	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	return names, nil
}

func (store *DirStore) LoadItem(name string) ([]byte, error) {
	fpath := filepath.Join(store.root, name)

	d, err := os.ReadFile(fpath)
	if err != nil {
		err = fmt.Errorf("failed to read item '%s': %w", name, err)
		return []byte{}, err
	}

	return d, nil
}

func (store *DirStore) StoreItem(name string, data []byte) error {
	fpath := filepath.Join(store.root, name)

	if err := os.WriteFile(fpath, data, 0644); err != nil {
		return fmt.Errorf("failed to read item '%s': %w", name, err)
	}

	return nil
}
