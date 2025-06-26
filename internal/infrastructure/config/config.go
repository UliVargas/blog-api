package config

import "os"

type Config struct {
	DBDSN     string
	JWTSECRET string
	PORT      string
}

func Load() *Config {
	return &Config{
		DBDSN:     os.Getenv("DBDSN"),
		JWTSECRET: os.Getenv("JWTSECRET"),
		PORT:      os.Getenv("PORT"),
	}
}