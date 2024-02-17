package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port     string `env:"PORT"`
	Postgres string `env:"Postgres"`
}

func GetConfig() (*Config, error) {
	config := Config{}
	err := cleanenv.ReadConfig(".env", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
