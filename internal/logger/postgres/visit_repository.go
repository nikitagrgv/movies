package postgres

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikitagrgv/movies/internal/logger"
)

type VisitRepository struct {
	pool *pgxpool.Pool
}

func NewVisitRepository(pool *pgxpool.Pool) *VisitRepository {
	return &VisitRepository{pool: pool}
}

func (r *VisitRepository) Create(ctx context.Context, request logger.CreateVisitRequest) (logger.Visit, error) {
	query := `
INSERT INTO logs.visits (ip_address, path, duration, attempted_at)
VALUES ($1, $2, $3, $4)
RETURNING id
`

	var id int
	err := r.pool.QueryRow(ctx, query,
		request.IP.String(),
		request.Path,
		request.Duration.Nanoseconds(),
		request.AttemptedAt).Scan(&id)
	if err != nil {
		return logger.Visit{}, err
	}

	return logger.Visit{
		ID:          id,
		IP:          request.IP,
		Path:        request.Path,
		Duration:    request.Duration,
		AttemptedAt: request.AttemptedAt,
	}, nil
}

func (r *VisitRepository) GetVisits(ctx context.Context, limit, offset int) ([]logger.Visit, error) {
	query := `
SELECT id, ip_address, path, duration, attempted_at
FROM logs.visits
ORDER BY attempted_at DESC
LIMIT $1 OFFSET $2
`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visits []logger.Visit
	for rows.Next() {
		var visit logger.Visit
		var durationNs int64
		var ipStr string
		if err := rows.Scan(&visit.ID, &ipStr, &visit.Path, &durationNs, &visit.AttemptedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		visit.Duration = time.Duration(durationNs)
		visit.IP, err = netip.ParseAddr(ipStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing IP: %w", err)
		}

		visits = append(visits, visit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return visits, nil
}
