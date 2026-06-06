package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikitagrgv/movies/internal/logger"
)

type VisitRepository struct {
	pool *pgxpool.Pool
}

func NewVisitRepository(pool *pgxpool.Pool) *VisitRepository {
	return &VisitRepository{pool: pool}
}

func (v *VisitRepository) Create(request logger.CreateVisitRequest) (logger.Visit, error) {
	return logger.Visit{}, nil
}

func (v *VisitRepository) GetVisits(limit, offset int) ([]logger.Visit, error) {
	return []logger.Visit{}, nil
}
