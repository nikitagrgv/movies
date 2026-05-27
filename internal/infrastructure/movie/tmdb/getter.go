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
	media := domain.Media{
		ID:          raw.ID,
		Title:       raw.Title,
		Overview:    raw.Overview,
		PosterURL:   poster,
		ReleaseYear: parseYear(raw.ReleaseDate),
	}
	res := domain.Movie{Media: media}

	return res, nil
}

func (g MovieGetter) GetTvShow(ctx context.Context, id int) (domain.TvShow, error) {
	var raw GetTvShowResponse
	err := g.client.get(
		ctx,
		"/tv/"+strconv.Itoa(id),
		nil,
		&raw,
	)
	if err != nil {
		return domain.TvShow{}, err
	}

	poster := g.client.getImageURL(raw.PosterPath)
	media := domain.Media{
		ID:          raw.ID,
		Title:       raw.Name,
		Overview:    raw.Overview,
		PosterURL:   poster,
		ReleaseYear: parseYear(raw.FirstAirDate),
	}

	// TODO#
	episodes := []domain.Episode{
		{EpisodeNumber: 1, SeasonNumber: 1, Name: "First Episode"},
	}
	seasons := []domain.Season{
		{SeasonNumber: 1, Name: "First Season", Episodes: episodes},
	}

	res := domain.TvShow{Media: media, Seasons: seasons}

	return res, nil
}
