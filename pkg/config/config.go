// Package config responsible for all the configuration to run the project.
package config

import "os"

// ApplicationEnvironment represent the application profile
type ApplicationEnvironment string

const (
	// DevelopmentEnvironment represent the application profile for development
	DevelopmentEnvironment ApplicationEnvironment = "dev"
	// ProductionEnvironment represent the application profile for production
	ProductionEnvironment ApplicationEnvironment = "prod"
	// QualityAssuranceAEnvironment represent the application profile for qa
	QualityAssuranceAEnvironment ApplicationEnvironment = "qa"
)

var config = map[string]string{
	"application":  "beeru",
	"environment":  getEnvVarWithDefault("APP_ENV", string(DevelopmentEnvironment)),
	"logLevel":     getEnvVarWithDefault("LOG_LEVEL", "debug"),
	"version":      getEnvVarWithDefault("VERSION", "detached"),
	"httpPort":     getEnvVarWithDefault("HTTP_PORT", ":8000"),
	"pgConnection": getEnvVarWithDefault("PG_CONNECTION", "user=caires dbname=postgres sslmode=disable"),
}

// Get config from config map
func Get(key string) string {
	return config[key]
}

func getEnvVarWithDefault(env string, defaultValue string) string {
	variable := defaultValue

	if setValue := os.Getenv(env); setValue != "" {
		variable = setValue
	}

	return variable
}
