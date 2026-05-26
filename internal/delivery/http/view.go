package http

import "github.com/nikitagrgv/movies/internal/domain"

type ErrorPageData struct {
	ErrorCode        int
	ErrorTitle       string
	ErrorDescription string
}

type SearchPageData struct {
	SearchString string
	MediaType    string
	CurrentPage  int
	TotalPages   int
	PrevPage     int
	NextPage     int
	Medias       []domain.MediaBase
}

type MovieView struct {
	ID          int
	Title       string
	Overview    string
	PosterURL   string
	ReleaseYear string
}

type TvShowView struct {
	ID          int
	Title       string
	Overview    string
	PosterURL   string
	ReleaseYear string
}
