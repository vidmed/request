package request

import (
	"sync"
	"time"

	"github.com/vidmed/logger"
)

// StoreService is a representation of request store entity.
// It`s responsibilities are
// - stores 50 store (generated on start)
// - with some periodicity it cancels one request and creates another one
// It`s safe for concurrent use.
type StoreService struct {
	letters []byte

	rsMux sync.RWMutex
	rs    [50]Request

	done chan struct{}
}

// NewStoreService creates StoreService with given update periodicity
// When StoreService is created - it generates 50 store
// and every p(time.Duration) it cancels one request and creates another one
func NewStoreService(p time.Duration) *StoreService {
	s := &StoreService{
		letters: []byte("abcdefghijklmnopqrstuvwxyz"),
		done:    make(chan struct{}),
	}

	for i := range s.rs {
		s.rs[i] = s.generateRequest()
	}

	go s.update(p)

	return s
}

// GetRequest reads a pseudo-random request from store and returns it
func (s *StoreService) GetRequest() Request {
	_, r := s.getRequest()
	return r
}

// Close closes StoreService
func (s *StoreService) Close() {
	logger.Get().Debug("StoreService: Close")

	if s.Closed() {
		return
	}
	close(s.done)
}

// Closed checks if StoreService is closed
func (s *StoreService) Closed() bool {
	select {
	case <-s.done:
		return true
	default:
		return false
	}
}

// getRequest reads a pseudo-random request from store and returns it and it`s index in store
func (s *StoreService) getRequest() (int, Request) {
	s.rsMux.RLock()
	defer s.rsMux.RUnlock()
	i := RandInt(len(s.rs))
	return i, s.rs[i]
}

// setRequest sets r Request to store at index i
func (s *StoreService) setRequest(i int, r Request) {
	s.rsMux.Lock()
	defer s.rsMux.Unlock()
	s.rs[i] = r
}

// generateRequest generates Request
func (s *StoreService) generateRequest() Request {
	b := Request{}
	for i := range b {
		b[i] = s.letters[RandInt(len(s.letters))]
	}
	return b
}

// replaceRequest gets pseudo-random index at store and replaces Request at that index
func (s *StoreService) replaceRequest() {
	i, r := s.getRequest()
	newR := s.generateRequest()
	s.setRequest(i, newR)

	logger.Get().Debugf("StoreService: request replaced: index - %d, old - %s, new - %s", i, r.String(), newR.String())
}

// update is blocking method. It takes periodicity argument (p)
// and call replaceRequest every p
func (s *StoreService) update(p time.Duration) {
	ticker := time.NewTicker(p)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			logger.Get().Debug("StoreService: update terminated")
			return
		case <-ticker.C:
			s.replaceRequest()
		}
	}
}
