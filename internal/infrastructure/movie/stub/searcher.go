package stub

import (
	"context"
	"errors"
	"strconv"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieSearcher struct{}

func NewMovieSearcher() *MovieSearcher {
	return &MovieSearcher{}
}

func (MovieSearcher) SearchMovies(ctx context.Context, query string, page int) (domain.SearchResult, error) {
	if page < 1 || page > 3 {
		return domain.SearchResult{}, errors.New("invalid page")
	}

	var movies []domain.Movie = make([]domain.Movie, 0, 3)
	if page == 1 {
		movies = append(movies, genMovie(query, 1))
		movies = append(movies, genMovie(query, 2))
		movies = append(movies, genMovie(query, 3))
	}
	if page == 2 {
		movies = append(movies, genMovie(query, 4))
		movies = append(movies, genMovie(query, 5))
		movies = append(movies, genMovie(query, 6))
	}
	if page == 3 {
		movies = append(movies, genMovie(query, 7))
		movies = append(movies, genMovie(query, 8))
		movies = append(movies, genMovie(query, 9))
	}

	return domain.SearchResult{Items: movies, CurrentPage: page, TotalPages: 3}, nil
}

func (s MovieSearcher) SearchTvShows(ctx context.Context, query string, page int) (domain.SearchResult, error) {
	//TODO implement me
	panic("implement me")
}

func genMovie(query string, index int) domain.Movie {
	name := query + " " + strconv.Itoa(index)
	media := domain.Media{
		ID:          index,
		Title:       name,
		Overview:    name + " is a beautiful movie about love... I cried!",
		PosterURL:   "",
		ReleaseYear: 2021,
	}
	movie := domain.Movie{media}
	return movie
}
