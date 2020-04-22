package request

type Service struct {
	store *StoreService
	views *ViewsService
}

func NewRequestService() *Service {
	s := &Service{
		store: NewStoreService(requestUpdatePeridicity),
		views: NewViewsService(),
	}

	return s
}

// GetRequest returns
func (s *Service) GetRequest() Request {
	r := s.store.GetRequest()
	s.views.IncreaseViews(r)

	return r
}

// GetViews returns number of requests` views
func (s *Service) GetViews() Views {
	return s.views.GetViews()
}

func (s *Service) Close() {
	s.store.Close()
}
