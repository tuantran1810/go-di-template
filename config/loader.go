package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

func LoadConfig[T any]() (T, error) {
	var conf T
	if err := env.Parse(&conf); err != nil {
		return conf, fmt.Errorf("failed to parse server config: %w", err)
	}

	return conf, nil
}

func MustLoadConfig[T any]() T {
	conf, err := LoadConfig[T]()
	if err != nil {
		panic(err)
	}

	return conf
}
