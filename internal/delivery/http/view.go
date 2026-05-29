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

type WatchServerView struct {
	ID   string
	Name string
}

type MovieView struct {
	ID          int
	Title       string
	Overview    string
	PosterURL   string
	ReleaseYear int

	CurrentServer string
	Servers       []WatchServerView
}

type EpisodeView struct {
	EpisodeNumber int
	Name          string
	IsAvailable   bool
}

type SeasonView struct {
	SeasonNumber int
	Name         string
	EpisodeCount int
}

type TvShowView struct {
	ID          int
	Title       string
	Overview    string
	PosterURL   string
	ReleaseYear int

	CurrentSeason  int
	CurrentEpisode int
	CurrentServer  string

	WatchURL string

	Seasons  []SeasonView
	Episodes []EpisodeView
	Servers  []WatchServerView
}
