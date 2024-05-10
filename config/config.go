package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"io/fs"
)

type Config struct {
	Port                 int    `env:"PORT"`
	Postgres             string `env:"POSTGRES"`
	RedisAddr            string `env:"REDIS_ADDR"`
	KafkaAddr            string `env:"KAFKA_ADDR"`
	KafkaTopicCreateUser string `env:"KAFKA_TOPIC_CREATE_USER"`
	MailServiceAddr      string `env:"MAIL_SERVICE_ADDR"`
	AppURL               string `env:"APP_URL"`
	LogJson              bool   `env:"LOG_JSON"`
	IntervalTime         string `env:"INTERVAL_TIME"`
}

func GetConfig() (*Config, error) {
	config := Config{}
	err := cleanenv.ReadConfig(".env", &config)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err = cleanenv.ReadEnv(&config)
			if err != nil {
				return nil, err
			}
			return &config, nil
		}
		return nil, err
	}
	return &config, nil
}

func (c *Config) Validation() error {
	if c.Port == 0 {
		c.Port = 8021
	}
	if c.Postgres == "" {
		return errors.New("postgres is empty")
	}
	if c.RedisAddr == "" {
		return errors.New("redis address is empty")
	}
	if c.KafkaAddr == "" {
		return errors.New("kafka address is empty")
	}
	if c.KafkaTopicCreateUser == "" {
		return errors.New("kafka topic create user is empty")
	}
	if c.MailServiceAddr == "" {
		return errors.New("mail service address is empty")
	}
	return nil
}
