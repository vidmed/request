package request

import "sync"

// ViewsService is a service that sores requests views
// It`s safe for concurrent use
type ViewsService struct {
	mux sync.RWMutex
	vs  Views
}

// NewViewsService creates new ViewsService instance
func NewViewsService() *ViewsService {
	return &ViewsService{
		vs: make(Views),
	}
}

// IncreaseViews gets Request as a parameter and increases it`s views count
func (v *ViewsService) IncreaseViews(r Request) {
	v.mux.Lock()
	defer v.mux.Unlock()
	v.vs[r]++
}

// GetViews returns the copy of Views
func (v *ViewsService) GetViews() Views {
	v.mux.RLock()
	defer v.mux.RUnlock()

	cp := make(Views, len(v.vs))
	for index, element := range v.vs {
		cp[index] = element
	}

	return cp
}
