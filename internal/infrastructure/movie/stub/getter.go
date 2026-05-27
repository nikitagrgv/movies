package stub

import (
	"context"
	"strconv"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieGetter struct{}

func NewMovieGetter() *MovieGetter {
	return &MovieGetter{}
}

func (MovieGetter) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	name := "Movie " + strconv.Itoa(id)
	media := domain.Media{
		ID:          id,
		Title:       name,
		Overview:    name + " is a beautiful movie about love... I cried!",
		PosterURL:   "",
		ReleaseYear: 2021,
	}
	movie := domain.Movie{Media: media}
	return movie, nil
}

func (g MovieGetter) GetTvShow(ctx context.Context, id int) (domain.TvShow, error) {
	name := "Movie " + strconv.Itoa(id)
	media := domain.Media{
		ID:          id,
		Title:       name,
		Overview:    name + " is a beautiful movie about love... I cried!",
		PosterURL:   "",
		ReleaseYear: 2021,
	}
	episodes := []domain.Episode{
		{EpisodeNumber: 1, SeasonNumber: 1, Name: "First Episode"},
	}
	seasons := []domain.Season{
		{SeasonNumber: 1, Name: "First Season", Episodes: episodes},
	}

	movie := domain.TvShow{Media: media, Seasons: seasons}
	return movie, nil
}
