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

func (MovieSearcher) SearchMovie(ctx context.Context, query string, page int) (domain.SearchMovieResult, error) {
	if page < 1 || page > 3 {
		return domain.SearchMovieResult{}, errors.New("invalid page")
	}

	var movies []domain.Movie
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

	return domain.SearchMovieResult{Movies: movies, CurrentPage: page, TotalPages: 3}, nil
}

func genMovie(query string, index int) domain.Movie {
	name := query + " " + strconv.Itoa(index)
	movie := domain.Movie{
		ID:          index,
		Title:       "title " + name,
		Overview:    "overview " + name,
		PosterURL:   "",
		ReleaseDate: "01-01-2021",
	}
	return movie
}
