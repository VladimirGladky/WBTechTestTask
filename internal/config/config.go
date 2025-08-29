package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	Kafka struct {
		KafkaBrokers []string `yaml:"kafka_brokers" env-default:"localhost:9092"`
		KafkaTopic   string   `yaml:"kafka_topic" env-default:"orders"`
		KafkaGroupID string   `yaml:"kafka_group_id" env-default:"orders"`
	}
)

type Config struct {
	Postgres PG     `yaml:"Postgres"`
	Port     string `yaml:"port" env-default:"4141"`
	Host     string `yaml:"host" env-default:"localhost"`
	Kafka    Kafka  `yaml:"Kafka"`
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
