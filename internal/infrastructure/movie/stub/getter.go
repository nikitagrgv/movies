package stub

import (
	"context"
	"strconv"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieGetter struct{}

func NewMovieGetter() *MovieGetter {
	return &MovieGetter{}
}

func (MovieGetter) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	name := "Movie " + strconv.Itoa(id)
	base := domain.MediaBase{
		ID:          id,
		Title:       name,
		Overview:    name + " is a beautiful movie about love... I cried!",
		PosterURL:   "",
		ReleaseDate: "01-01-2021",
	}
	movie := domain.Movie{Base: base}
	return movie, nil
}

func (g MovieGetter) GetTvShow(ctx context.Context, id int) (domain.TvShow, error) {
	//TODO implement me
	panic("implement me")
}
