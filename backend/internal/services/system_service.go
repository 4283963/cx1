package services

import (
	"sync"
)

type SystemService struct {
	mu            sync.RWMutex
	forceMode     bool
	forceModeChan chan bool
}

var systemService *SystemService
var systemOnce sync.Once

func NewSystemService() *SystemService {
	systemOnce.Do(func() {
		systemService = &SystemService{
			forceModeChan: make(chan bool, 10),
		}
	})
	return systemService
}

func (s *SystemService) SetForceMode(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.forceMode = enabled
	select {
	case s.forceModeChan <- enabled:
	default:
	}
}

func (s *SystemService) GetForceMode() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.forceMode
}

func (s *SystemService) GetForceModeChan() <-chan bool {
	return s.forceModeChan
}
