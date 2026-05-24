package tmdb

import "net/http"

type TmdbClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}
