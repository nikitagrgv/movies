package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	ListenPort int
	TmdbAPIKey string
	Stubs      []stubType
}

func LoadFromEnv() (Config, error) {
	portStr, ok := os.LookupEnv("MOVIES_LISTEN_PORT")
	if !ok {
		return Config{}, errors.New("MOVIES_LISTEN_PORT must be set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return Config{}, errors.New("MOVIES_LISTEN_PORT must be an integer")
	}

	tmdbKey, ok := os.LookupEnv("MOVIES_TMDB_KEY")
	if !ok {
		return Config{}, errors.New("MOVIES_TMDB_KEY must be set")
	}

	var stubs []stubType
	stubsStr, ok := os.LookupEnv("MOVIES_STUBS")
	if ok {
		s, err := parseStubTypes(stubsStr)
		if err != nil {
			return Config{}, err
		}
		stubs = s
	}

	return Config{
		ListenPort: port,
		TmdbAPIKey: tmdbKey,
		Stubs:      stubs,
	}, nil
}

func (c *Config) IsStubUsed(stub stubType) bool {
	if len(c.Stubs) == 0 {
		return false
	}

	for _, s := range c.Stubs {
		if s == stub {
			return true
		}
	}
	return false
}
