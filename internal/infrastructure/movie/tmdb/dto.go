package tmdb

type SearchMovieResponse struct {
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
	Results    []struct {
		ID           int    `json:"id"`
		Title        string `json:"title"`
		Overview     string `json:"overview"`
		FirstAirDate string `json:"first_air_date"`
		PosterPath   string `json:"poster_path"`
	} `json:"results"`
}

type GetMovieResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
	PosterPath  string `json:"poster_path"`
}
