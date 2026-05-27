package domain

import (
	"errors"
)

var ErrInvalidMediaType = errors.New("invalid media type")

type MediaType string

const (
	MovieType  MediaType = "movie"
	TvShowType MediaType = "tv"
)

func ParseMediaType(s string) (MediaType, error) {
	switch MediaType(s) {
	case MovieType, TvShowType:
		return MediaType(s), nil
	default:
		return "", ErrInvalidMediaType
	}
}

type MediaBase struct {
	ID          int
	Title       string
	Overview    string
	PosterURL   string
	ReleaseYear int
}

type Movie struct {
	MediaBase
}

type Episode struct {
	EpisodeNumber int
	SeasonNumber  int
	Name          string
}

type Season struct {
	SeasonNumber int
	Name         string
	Episodes     []Episode
}

type TvShow struct {
	MediaBase
	Seasons []Season
}
