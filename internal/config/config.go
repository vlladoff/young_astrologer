package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type Config struct {
	StorageDataSource string `env:"YA_STORAGE_DATA_SOURCE" env-required:"true"`
	APODEndpoint      string `env:"YA_APOD_ENDPOINT" env-required:"true"`
	HTTPServer
}

type HTTPServer struct {
	Address     string        `env:"YA_HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `env:"YA_HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"YA_HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
