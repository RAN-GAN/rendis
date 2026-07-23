package store

import (
	"encoding/json"
	"log"
	"time"
)

type Snapshot struct {
	Version   int
	CreatedAt time.Time
	Data      map[string]Entry
}

func (s *Store) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snapshot := Snapshot{
		Version:   1,
		CreatedAt: time.Now(),
		Data:      s.data,
	}
	return json.Marshal(snapshot)
}

func (s *Store) Restore(data []byte) error {
	var snapshot Snapshot

	if err := json.Unmarshal(data, &snapshot); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = snapshot.Data
	log.Println("Restored data from:", snapshot.CreatedAt)

	return nil
}
