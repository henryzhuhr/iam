package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	DSN string
}

func LoadDBConfig() (DBConfig, error) {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		return DBConfig{}, fmt.Errorf("PG_DSN is required")
	}
	return DBConfig{DSN: dsn}, nil
}
