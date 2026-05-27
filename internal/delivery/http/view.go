package http

type ErrorPageView struct {
	ErrorCode        int
	ErrorTitle       string
	ErrorDescription string
}

type SearchItemView struct {
	ID          int
	Title       string
	Overview    string
	PosterURL   string
	ReleaseYear int
}

type SearchView struct {
	SearchString string
	MediaType    string
	CurrentPage  int
	TotalPages   int
	PrevPage     int
	NextPage     int
	Items        []SearchItemView
}

type MovieView struct {
	ID          int
	Title       string
	Overview    string
	PosterURL   string
	ReleaseYear int
}

type TvShowView struct {
	ID          int
	Title       string
	Overview    string
	PosterURL   string
	ReleaseYear int
}
