package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	configs "github.com/nikitagrgv/movies/config"
	"gopkg.in/yaml.v3"
)

type DbConfig struct {
	User     string
	Password string
	Host     string
	Db       string
	Port     int
}

type Config struct {
	ListenPort   int
	TmdbToken    string
	Stubs        []stubType
	WatchServers []WatchServerConfig

	GRPCListenPort int

	Db DbConfig
}

type WatchServersConfig struct {
	Servers []WatchServerConfig `yaml:"servers"`
}

type WatchServerConfig struct {
	ID                string `yaml:"id"`
	Name              string `yaml:"name"`
	MovieURLTemplate  string `yaml:"movie_url_template"`
	TvShowURLTemplate string `yaml:"tv_show_url_template"`
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

	db, err := loadDbConfig()
	if err != nil {
		return Config{}, err
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

	servers, err := loadWatchServersConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		ListenPort:   port,
		TmdbToken:    tmdbToken,
		Stubs:        stubs,
		WatchServers: servers,
		Db:           db,
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

func loadWatchServersConfig() ([]WatchServerConfig, error) {
	var cfg WatchServersConfig
	err := yaml.Unmarshal(configs.WatchServersRawConfig, &cfg)
	if err != nil {
		return []WatchServerConfig{}, fmt.Errorf("error parsing watch servers config: %s", err)
	}
	return cfg.Servers, nil
}

func loadDbConfig() (DbConfig, error) {
	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return DbConfig{}, errors.New("POSTGRES_USER must be set")
	}

	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return DbConfig{}, errors.New("POSTGRES_PASSWORD must be set")
	}

	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		return DbConfig{}, errors.New("POSTGRES_HOST must be set")
	}

	db, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		return DbConfig{}, errors.New("POSTGRES_DB must be set")
	}

	return DbConfig{
		User:     user,
		Password: password,
		Host:     host,
		Db:       db,
		Port:     5432, // TODO: #hardcoded
	}, nil
}
