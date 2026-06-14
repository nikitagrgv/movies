package tmdb

import (
	"context"
	"strconv"
	"time"

	"github.com/nikitagrgv/movies/internal/media"
)

type MediaGetter struct {
	client *Client
}

func NewMediaGetter(client *Client) *MediaGetter {
	return &MediaGetter{client: client}
}

func (g *MediaGetter) GetMovie(ctx context.Context, id int) (media.Movie, error) {
	var raw GetMovieResponse
	err := g.client.get(
		ctx,
		"/movie/"+strconv.Itoa(id),
		nil,
		&raw,
	)
	if err != nil {
		return media.Movie{}, err
	}

	poster := g.client.getImageURL(raw.PosterPath)
	m := media.Media{
		ID:          raw.ID,
		Title:       raw.Title,
		Overview:    raw.Overview,
		PosterURL:   poster,
		ReleaseYear: parseYear(raw.ReleaseDate),
	}
	res := media.Movie{Media: m}

	return res, nil
}

func (g *MediaGetter) GetTvShow(ctx context.Context, id int) (media.TvShow, error) {
	var raw GetTvShowResponse
	err := g.client.get(
		ctx,
		"/tv/"+strconv.Itoa(id),
		nil,
		&raw,
	)
	if err != nil {
		return media.TvShow{}, err
	}

	poster := g.client.getImageURL(raw.PosterPath)
	m := media.Media{
		ID:          raw.ID,
		Title:       raw.Name,
		Overview:    raw.Overview,
		PosterURL:   poster,
		ReleaseYear: parseYear(raw.FirstAirDate),
	}

	res := media.TvShow{Media: m, TotalSeasons: raw.NumSeasons}

	return res, nil
}

func (g *MediaGetter) GetTvShowSeason(ctx context.Context, id int, seasonNumber int) (media.Season, error) {
	var raw GetSeasonResponse
	err := g.client.get(
		ctx,
		"/tv/"+strconv.Itoa(id)+"/season/"+strconv.Itoa(seasonNumber),
		nil,
		&raw,
	)
	if err != nil {
		return media.Season{}, err
	}

	var episodes []media.Episode
	for _, rawEpisode := range raw.Episodes {
		date, err := time.Parse("2006-01-02", rawEpisode.AirDate)
		if err != nil {
			date = time.Time{}
		}
		episode := media.Episode{
			EpisodeNumber: rawEpisode.EpisodeNumber,
			SeasonNumber:  seasonNumber,
			Name:          rawEpisode.Name,
			Date:          date,
		}
		episodes = append(episodes, episode)
	}

	season := media.Season{
		ShowID:       id,
		SeasonNumber: seasonNumber,
		Name:         raw.Name,
		Episodes:     episodes,
	}
	return season, nil
}
