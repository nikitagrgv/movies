package http

import "github.com/nikitagrgv/movies/internal/domain"

type ErrorPageData struct {
	ErrorCode        int
	ErrorDescription string
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
	MovieID int
	Title   string
}
