package domain

import "errors"

type MediaType string

const (
	MovieType  MediaType = "movie"
	TvShowType MediaType = "tv"
)

func ParseMediaType(s string) (MediaType, error) {
	switch s {
	case "movie":
		return MovieType, nil
	case "tv":
		return TvShowType, nil
	default:
		return "", errors.New("invalid media type")
	}
}

type MediaBase struct {
	ID          int
	Title       string
	Overview    string
	PosterURL   string
	ReleaseDate string
}

type Movie struct {
	Base MediaBase
}

type TvShow struct {
	Base MediaBase
}
