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

type Media struct {
	ID          int
	Title       string
	Overview    string
	PosterURL   string
	ReleaseYear int
}

type Movie struct {
	Media
}

type TvShow struct {
	Media
	TotalSeasons int
}

type Episode struct {
	EpisodeNumber int
	SeasonNumber  int
	Name          string
}

type Season struct {
	ShowID       int
	SeasonNumber int
	Name         string
	Episodes     []Episode
}

type WatchServer struct {
	ID   string
	Name string
}
