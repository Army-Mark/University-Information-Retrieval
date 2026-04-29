package config

import (
	"os"
)

type Config struct {
	Debug  bool
	DBPath string
	Port   string
	Host   string
}

func Load() *Config {
	debug := false
	if os.Getenv("FLASK_DEBUG") == "true" {
		debug = true
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "school.db"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	return &Config{
		Debug:  debug,
		DBPath: dbPath,
		Port:   port,
		Host:   host,
	}
}
