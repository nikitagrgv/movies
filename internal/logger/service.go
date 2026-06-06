package logger

type VisitRepository interface {
	Create(CreateVisitRequest) (Visit, error)
	GetVisits(limit, offset int) ([]Visit, error)
}

type Service struct {
	repo VisitRepository
}

func NewService(repo VisitRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateVisit(req CreateVisitRequest) (Visit, error) {
	return s.repo.Create(req)
}

func (s *Service) GetVisits(limit, offset int) ([]Visit, error) {
	return s.repo.GetVisits(limit, offset)
}
