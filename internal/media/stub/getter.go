package stub

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/nikitagrgv/movies/internal/media"
)

type MediaGetter struct{}

func NewMediaGetter() *MediaGetter {
	return &MediaGetter{}
}

func (MediaGetter) GetMovie(ctx context.Context, id int) (media.Movie, error) {
	name := "Movie " + strconv.Itoa(id)
	m := media.Media{
		ID:          id,
		Title:       name,
		Overview:    name + " is a beautiful movie about love... I cried!",
		PosterURL:   "",
		ReleaseYear: 2021,
	}
	movie := media.Movie{Media: m}
	return movie, nil
}

func (g MediaGetter) GetTvShow(ctx context.Context, id int) (media.TvShow, error) {
	name := "Movie " + strconv.Itoa(id)
	m := media.Media{
		ID:          id,
		Title:       name,
		Overview:    name + " is a beautiful movie about love... I cried!",
		PosterURL:   "",
		ReleaseYear: 2021,
	}

	movie := media.TvShow{Media: m, TotalSeasons: 2}
	return movie, nil
}

func (g MediaGetter) GetTvShowSeason(ctx context.Context, id, season int) (media.Season, error) {
	if season < 0 || season > 2 {
		return media.Season{}, errors.New("season number must be between 0 and 1")
	}

	date, _ := time.Parse("01 02 2006", "01 02 2006")
	episodes := []media.Episode{
		{EpisodeNumber: 1, SeasonNumber: season, Name: "First Episode", Date: date},
		{EpisodeNumber: 2, SeasonNumber: season, Name: "Second Episode", Date: date.Add(time.Hour * 24 * 10)},
		{EpisodeNumber: 3, SeasonNumber: season, Name: "Third Episode", Date: date.Add(time.Hour * 24 * 10)},
	}

	return media.Season{
		ShowID:       id,
		SeasonNumber: 1,
		Name:         "Season " + strconv.Itoa(season),
		Episodes:     episodes,
	}, nil
}
