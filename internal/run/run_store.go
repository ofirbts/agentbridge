package run

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type RunStore struct {
	dir string
	mu  sync.RWMutex
	mem map[string]*Run
}

func NewRunStore(dir string) *RunStore {
	if dir == "" {
		dir = ".agentbridge/runs"
	}
	return &RunStore{
		dir: dir,
		mem: map[string]*Run{},
	}
}

func (s *RunStore) Insert(id string, run *Run) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	copied := *run
	copied.ID = id
	s.mem[id] = &copied
	return s.persist(id, &copied)
}

func (s *RunStore) Get(id string) (*Run, error) {
	s.mu.RLock()
	if r, ok := s.mem[id]; ok {
		s.mu.RUnlock()
		copied := *r
		return &copied, nil
	}
	s.mu.RUnlock()

	path := filepath.Join(s.dir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("run %q not found: %w", id, err)
	}
	return decodeRun(data)
}

func decodeRun(data []byte) (*Run, error) {
	var aux struct {
		Run
		LegacyLogs []StepLog `json:"logs"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, err
	}
	r := aux.Run
	if len(r.StepLogs) == 0 && len(aux.LegacyLogs) > 0 {
		r.StepLogs = aux.LegacyLogs
	}
	return &r, nil
}

func (s *RunStore) persist(id string, run *Run) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.dir, id+".json"), data, 0o644)
}
