package collection

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
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
		return nil, &ErrNotDir{Path: root}
	}

	store := &DirStore{root: root}

	return store, nil
}

func (store *DirStore) MakeSubstore(name string) error {
	path := filepath.Join(store.root, name)

	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to make substore '%s': %w", name, err)
	}

	return nil
}

func (store *DirStore) SubstoreNames() []string {
	entries, err := os.ReadDir(store.root)
	if err != nil {
		log.Error().Err(err).Str("dir", store.root).Msg("failed to read store dir")

		return []string{}
	}

	names := []string{}

	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}

	return names
}

func (store *DirStore) ItemNames() []string {
	entries, err := os.ReadDir(store.root)
	if err != nil {
		log.Error().Err(err).Str("dir", store.root).Msg("failed to read store dir")

		return []string{}
	}

	names := []string{}

	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}

	return names
}

func (store *DirStore) Item(name string) (Item, error) {
	fpath := filepath.Join(store.root, name)

	return NewFileItem(fpath)
}

func (store *DirStore) Substore(name string) (Store, error) {
	sub := filepath.Join(store.root, name)

	return NewDirStore(sub)
}
