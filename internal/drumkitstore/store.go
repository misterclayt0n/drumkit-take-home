package drumkitstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"drumkit-take-home/internal/load"
)

type Store struct {
	path string

	mu      sync.RWMutex
	records map[string]load.Load
}

type persistedStore struct {
	Records map[string]load.Load `json:"records"`
}

func New(path string) (*Store, error) {
	store := &Store{
		path:    path,
		records: map[string]load.Load{},
	}

	if err := store.loadFromDisk(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) Save(turvoID string, payload load.Load) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records[turvoID] = payload
	return s.saveToDiskLocked()
}

func (s *Store) MergeList(loads []load.Load) []load.Load {
	merged := make([]load.Load, 0, len(loads))
	for _, item := range loads {
		merged = append(merged, s.Merge(item))
	}
	return merged
}

func (s *Store) Merge(item load.Load) load.Load {
	s.mu.RLock()
	defer s.mu.RUnlock()

	saved, ok := s.records[item.ExternalTMSLoadID]
	if !ok {
		return item
	}

	merged := saved
	if item.Status != "" {
		merged.Status = item.Status
	}
	return merged
}

func (s *Store) loadFromDisk() error {
	if s.path == "" {
		return nil
	}

	content, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read drumkit store: %w", err)
	}
	if len(content) == 0 {
		return nil
	}

	var persisted persistedStore
	if err := json.Unmarshal(content, &persisted); err != nil {
		return fmt.Errorf("decode drumkit store: %w", err)
	}
	if persisted.Records != nil {
		s.records = persisted.Records
	}

	return nil
}

func (s *Store) saveToDiskLocked() error {
	if s.path == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create drumkit store directory: %w", err)
	}

	content, err := json.MarshalIndent(persistedStore{Records: s.records}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode drumkit store: %w", err)
	}

	tempPath := s.path + ".tmp"
	if err := os.WriteFile(tempPath, content, 0o644); err != nil {
		return fmt.Errorf("write drumkit store temp file: %w", err)
	}
	if err := os.Rename(tempPath, s.path); err != nil {
		return fmt.Errorf("replace drumkit store file: %w", err)
	}

	return nil
}
