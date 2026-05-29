package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	configs "github.com/nikitagrgv/movies/config"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenPort    int
	TmdbToken     string
	Stubs         []stubType
	WatchServices []WatchServiceConfig
}

type WatchServicesConfig struct {
	Services []WatchServiceConfig `yaml:"services"`
}

type WatchServiceConfig struct {
	ID                int `yaml:"id"`
	Name              int `yaml:"name"`
	MovieURLTemplate  int `yaml:"movie_url_template"`
	TvShowURLTemplate int `yaml:"tv_show_url_template"`
}

func Load() (Config, error) {
	portStr, ok := os.LookupEnv("MOVIES_LISTEN_PORT")
	if !ok {
		return Config{}, errors.New("MOVIES_LISTEN_PORT must be set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return Config{}, errors.New("MOVIES_LISTEN_PORT must be an integer")
	}

	tmdbToken, ok := os.LookupEnv("MOVIES_TMDB_TOKEN")
	if !ok {
		return Config{}, errors.New("MOVIES_TMDB_TOKEN must be set")
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

	services, err := loadWatchServicesConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		ListenPort:    port,
		TmdbToken:     tmdbToken,
		Stubs:         stubs,
		WatchServices: services,
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

func loadWatchServicesConfig() ([]WatchServiceConfig, error) {
	var cfg WatchServicesConfig
	err := yaml.Unmarshal(configs.WatchServicesRawConfig, &cfg)
	if err != nil {
		return []WatchServiceConfig{}, fmt.Errorf("error parsing watch services config: %s", err)
	}
	return cfg.Services, nil
}
