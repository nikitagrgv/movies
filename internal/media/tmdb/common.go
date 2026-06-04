package tmdb

import "time"

func parseYear(releaseDate string) int {
	t, err := time.Parse("2006-01-02", releaseDate)
	if err != nil {
		return 0
	}
	return t.Year()
}
