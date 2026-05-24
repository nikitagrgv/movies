package usecase

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type GetMovieUsecase struct {
	getter     domain.MovieGetter
	noImageURL string
}

func NewGetMovieUsecase(getter domain.MovieGetter, noImageURL string) *GetMovieUsecase {
	return &GetMovieUsecase{getter: getter, noImageURL: noImageURL}
}

func (u *GetMovieUsecase) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	movie, err := u.getter.GetMovie(ctx, id)
	if err != nil {
		return domain.Movie{}, err
	}

	movie = u.normalizeMovie(movie)
	return movie, nil
}

func (u *GetMovieUsecase) normalizeMovie(movie domain.Movie) domain.Movie {
	if movie.PosterURL == "" {
		movie.PosterURL = u.noImageURL
	}

	return movie
}
