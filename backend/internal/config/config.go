package config

import (
	"os"
	"strings"
)

type Config struct {
	Port         string
	Host         string
	BaseURL      string
	AllowedHosts []string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		Host:         getEnv("HOST", "localhost"),
		BaseURL:      strings.TrimRight(getEnv("BASE_URL", ""), "/"),
		AllowedHosts: parseCSV(getEnv("CORS_ALLOW_ORIGINS", "*")),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}
