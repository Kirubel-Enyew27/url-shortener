package config

import "os"

type Config struct {
	Port string
	Host string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		Host: getEnv("HOST", "localhost"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
