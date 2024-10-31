package config

import "os"

// Environment values
const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

// GetEnvironment returns the current environment, defaulting to production
func GetEnvironment() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = EnvProduction
	}
	return env
}

// IsDevelopment returns true if running in development environment
func IsDevelopment() bool {
	return GetEnvironment() == EnvDevelopment
}
