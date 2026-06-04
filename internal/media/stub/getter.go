package stub

import (
	"context"
	"errors"
	"strconv"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MediaGetter struct{}

func NewMediaGetter() *MediaGetter {
	return &MediaGetter{}
}

func (MediaGetter) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
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

func (g MediaGetter) GetTvShow(ctx context.Context, id int) (domain.TvShow, error) {
	name := "Movie " + strconv.Itoa(id)
	media := domain.Media{
		ID:          id,
		Title:       name,
		Overview:    name + " is a beautiful movie about love... I cried!",
		PosterURL:   "",
		ReleaseYear: 2021,
	}

	movie := domain.TvShow{Media: media, TotalSeasons: 2}
	return movie, nil
}

func (g MediaGetter) GetTvShowSeason(ctx context.Context, id, season int) (domain.Season, error) {
	if season < 0 || season > 2 {
		return domain.Season{}, errors.New("season number must be between 0 and 1")
	}

	episodes := []domain.Episode{
		{EpisodeNumber: 1, SeasonNumber: season, Name: "First Episode"},
		{EpisodeNumber: 2, SeasonNumber: season, Name: "Second Episode"},
		{EpisodeNumber: 3, SeasonNumber: season, Name: "Third Episode"},
	}

	return domain.Season{
		ShowID:       id,
		SeasonNumber: 1,
		Name:         "Season " + strconv.Itoa(season),
		Episodes:     episodes,
	}, nil
}
