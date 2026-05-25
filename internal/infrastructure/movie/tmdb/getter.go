package tmdb

import (
	"context"
	"strconv"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieGetter struct {
	client *Client
}

func NewMovieGetter(client *Client) *MovieGetter {
	return &MovieGetter{client: client}
}

func (g MovieGetter) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	var raw GetMovieResponse
	err := g.client.get(
		ctx,
		"/movie/"+strconv.Itoa(id),
		nil,
		&raw,
	)
	if err != nil {
		return domain.Movie{}, err
	}

	poster := g.client.getImageURL(raw.PosterPath)
	base := domain.MediaBase{
		ID:          raw.ID,
		Title:       raw.Title,
		Overview:    raw.Overview,
		PosterURL:   poster,
		ReleaseDate: raw.ReleaseDate,
	}
	res := domain.Movie{Base: base}

	return res, nil
}

func (g MovieGetter) GetTvShow(ctx context.Context, id int) (domain.TvShow, error) {
	//TODO implement me
	panic("implement me")
}
