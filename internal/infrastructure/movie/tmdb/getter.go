package tmdb

import (
	"context"
	"strconv"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MediaGetter struct {
	client *Client
}

func NewMediaGetter(client *Client) *MediaGetter {
	return &MediaGetter{client: client}
}

func (g *MediaGetter) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
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

func (g *MediaGetter) GetTvShow(ctx context.Context, id int) (domain.TvShow, error) {
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

	res := domain.TvShow{Media: media, TotalSeasons: raw.NumSeasons}

	return res, nil
}

func (g *MediaGetter) GetTvShowSeason(ctx context.Context, id int, seasonNumber int) (domain.Season, error) {
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
		ShowID:       id,
		SeasonNumber: seasonNumber,
		Name:         raw.Name,
		Episodes:     episodes,
	}
	return season, nil
}
