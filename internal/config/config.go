package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	ListenPort int
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

	return Config{
		ListenPort: port,
	}, nil
}
