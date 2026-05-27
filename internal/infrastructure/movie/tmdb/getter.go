package tmdb

import (
	"context"
	"strconv"
	"sync"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieGetter struct {
	client *Client
}

func NewMovieGetter(client *Client) *MovieGetter {
	return &MovieGetter{client: client}
}

func (g *MovieGetter) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
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

func (g *MovieGetter) GetTvShow(ctx context.Context, id int) (domain.TvShow, error) {
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

	var seasons = make([]domain.Season, raw.NumSeasons)
	var wg sync.WaitGroup
	wg.Add(raw.NumSeasons)
	for seasonNumber := 1; seasonNumber <= raw.NumSeasons; seasonNumber++ {
		go func(id int, s int) {
			defer wg.Done()
			season, err := g.getSeason(ctx, id, s)
			if err != nil {
				season = domain.Season{
					SeasonNumber: s,
					Name:         "Invalid Season",
				}
			}
			seasons[seasonNumber-1] = season
		}(id, seasonNumber)
	}

	wg.Wait()

	res := domain.TvShow{Media: media, Seasons: seasons}

	return res, nil
}

func (g *MovieGetter) getSeason(ctx context.Context, id int, seasonNumber int) (domain.Season, error) {
	var raw GetSeasonResponse
	err := g.client.get(
		ctx,
		"/tv/"+strconv.Itoa(id)+"/season/"+strconv.Itoa(seasonNumber),
		nil,
		&raw,
	)
	if err != nil {
		return domain.Season{}, err
	}

	var episodes []domain.Episode
	for _, rawEpisode := range raw.Episodes {
		episode := domain.Episode{
			EpisodeNumber: rawEpisode.EpisodeNumber,
			SeasonNumber:  seasonNumber,
			Name:          rawEpisode.Name,
		}
		episodes = append(episodes, episode)
	}

	season := domain.Season{
		SeasonNumber: seasonNumber,
		Name:         raw.Name,
		Episodes:     episodes,
	}
	return season, nil
}
