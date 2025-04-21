package app

import (
	"os"
)

type Config struct {
	HTTPPort    string
	GRPCPort    string
	MetricsPort string
	DatabaseURL string
	JWTSecret   string
}

func LoadConfig() *Config {
	return &Config{
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		GRPCPort:    getEnv("GRPC_PORT", "3000"),
		MetricsPort: getEnv("METRICS_PORT", "9000"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/pvz?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}