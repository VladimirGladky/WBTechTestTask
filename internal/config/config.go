package config

import (
	"WBTechTestTask/pkg/postgres"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Postgres postgres.Config `yaml:"Postgres"`
	Port     string          `yaml:"port" env-default:"4047"`
	Host     string          `yaml:"host" env-default:"0.0.0.0"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	var cfg Config
	err := cleanenv.ReadConfig("./config/config.yaml", &cfg)
	if err != nil {
		return nil, err
	}
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
