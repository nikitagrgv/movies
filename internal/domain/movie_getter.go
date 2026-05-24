package domain

import "context"

type MovieGetter interface {
	GetMovie(ctx context.Context, id int) (Movie, error)
}
