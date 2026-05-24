package http

import "github.com/nikitagrgv/movies/internal/domain"

type ErrorPageData struct {
	ErrorCode  int
	ErrorTitle string
}

type SearchPageData struct {
	SearchString string
	CurrentPage  int
	TotalPages   int
	PrevPage     int
	NextPage     int
	Movies       []domain.Movie
}

type MovieView struct {
	MovieID     int
	Title       string
	Overview    string
	PosterURL   string
	ReleaseDate string
}
