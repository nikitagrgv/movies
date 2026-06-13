package logger

import "context"

type VisitRepository interface {
	Create(ctx context.Context, req CreateVisitRequest) (Visit, error)
	GetVisits(ctx context.Context, limit, offset int) ([]Visit, error)
}

type Service struct {
	repo VisitRepository
}

func NewService(repo VisitRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) PushVisit(req CreateVisitRequest) {
	// TODO# to workers queue
	_, _ = s.repo.Create(context.TODO(), req)
}

func (s *Service) GetVisits(ctx context.Context, limit, offset int) ([]Visit, error) {
	return s.repo.GetVisits(ctx, limit, offset)
}
