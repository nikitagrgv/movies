package stub

import (
	"context"
	"errors"
	"strconv"

	"github.com/nikitagrgv/movies/internal/media"
)

type MediaSearcher struct{}

func NewMediaSearcher() *MediaSearcher {
	return &MediaSearcher{}
}

func (MediaSearcher) SearchMovies(_ context.Context, query string, page int) (media.SearchResult, error) {
	return getMedias(query, page)
}

func (s MediaSearcher) SearchTvShows(_ context.Context, query string, page int) (media.SearchResult, error) {
	return getMedias(query, page)
}

func getMedias(query string, page int) (media.SearchResult, error) {
	if page < 1 || page > 3 {
		return media.SearchResult{}, errors.New("invalid page")
	}

	var movies = make([]media.Media, 0, 3)
	if page == 1 {
		movies = append(movies, getMedia(query, 1))
		movies = append(movies, getMedia(query, 2))
		movies = append(movies, getMedia(query, 3))
	}
	if page == 2 {
		movies = append(movies, getMedia(query, 4))
		movies = append(movies, getMedia(query, 5))
		movies = append(movies, getMedia(query, 6))
	}
	if page == 3 {
		movies = append(movies, getMedia(query, 7))
		movies = append(movies, getMedia(query, 8))
		movies = append(movies, getMedia(query, 9))
	}

	return media.SearchResult{Items: movies, CurrentPage: page, TotalPages: 3}, nil
}

func getMedia(query string, index int) media.Media {
	name := query + " " + strconv.Itoa(index)
	media := media.Media{
		ID:          index,
		Title:       name,
		Overview:    name + " is a beautiful movie about love... I cried!",
		PosterURL:   "",
		ReleaseYear: 2021,
	}
	return media
}
