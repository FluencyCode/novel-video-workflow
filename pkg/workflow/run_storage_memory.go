package workflow

import (
	"fmt"
	"sync"
)

// MemoryRunStorage is an in-memory implementation for testing.
type MemoryRunStorage struct {
	mu   sync.RWMutex
	runs map[uint]WorkflowRun
}

// NewMemoryRunStorage creates a new in-memory storage.
func NewMemoryRunStorage() *MemoryRunStorage {
	return &MemoryRunStorage{
		runs: map[uint]WorkflowRun{},
	}
}

// Save stores a workflow run.
func (s *MemoryRunStorage) Save(run WorkflowRun) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.runs[run.ChapterID] = run
	return nil
}

// Load retrieves a workflow run by chapter ID.
func (s *MemoryRunStorage) Load(chapterID uint) (WorkflowRun, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	run, ok := s.runs[chapterID]
	if !ok {
		return WorkflowRun{}, fmt.Errorf("run not found for chapter %d", chapterID)
	}
	return run, nil
}

// LoadByChapterID implements RunStorage interface (alias for Load).
func (s *MemoryRunStorage) LoadByChapterID(chapterID uint) (WorkflowRun, error) {
	return s.Load(chapterID)
}