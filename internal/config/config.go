package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Proxy     string  `env:"PROXY"`
	StoreID   uint64  `env:"ID_STORE" envDefault:"13151"`
	Latitude  float64 `env:"LATITUDE"`
	Longitude float64 `env:"LONGITUDE"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return cfg, nil
}

func (c *Config) Validate() error {

	if (c.Latitude != 0 && c.Longitude == 0) || (c.Latitude == 0 && c.Longitude != 0) {
		return fmt.Errorf("both LATITUDE and LONGITUDE must be set together")
	}

	return nil
}
