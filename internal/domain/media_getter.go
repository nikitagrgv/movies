package domain

import "context"

type MediaGetter interface {
	GetMovie(ctx context.Context, id int) (Movie, error)
	GetTvShow(ctx context.Context, id int) (TvShow, error)
}
