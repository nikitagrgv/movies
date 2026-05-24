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
	Movies       []domain.Movie
}
