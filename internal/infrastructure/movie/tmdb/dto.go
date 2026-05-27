package tmdb

type SearchMovieResponse struct {
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
	Results    []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Overview    string `json:"overview"`
		ReleaseDate string `json:"release_date"`
		PosterPath  string `json:"poster_path"`
	} `json:"results"`
}

type GetMovieResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
	PosterPath  string `json:"poster_path"`
}

type SearchTvShowResponse struct {
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
	Results    []struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Overview     string `json:"overview"`
		FirstAirDate string `json:"first_air_date"`
		PosterPath   string `json:"poster_path"`
	} `json:"results"`
}

type GetTvShowResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	FirstAirDate string `json:"first_air_date"`
	PosterPath   string `json:"poster_path"`
	Seasons      []struct {
		SeasonNumber int `json:"season_number"`
	} `json:"seasons"`
}

type GetSeasonResponse struct {
	Name     string `json:"name"`
	Episodes []struct {
		Name          string `json:"name"`
		EpisodeNumber int    `json:"episode_number"`
	} `json:"episodes"`
}
