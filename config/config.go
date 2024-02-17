package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Port string `json:"PORT"`
	Data string `json:"DATA"`
}

var config Config

func GetConfig() *Config {
	config = Config{}
	err := cleanenv.ReadConfig(".env", &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}
