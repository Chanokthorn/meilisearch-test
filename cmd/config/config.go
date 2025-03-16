package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Meilisearch
}

type Meilisearch struct {
	Host string `env:"MEILISEARCH_HOST" envDefault:"http://localhost:7700"`
	MasterKey string `env:"MEILISEARCH_MASTER_KEY" envDefault:"MASTER_KEY"`
}

func Read() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}
	
	return cfg
}