package config

import (
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Load() (*viper.Viper, error) {
	// Load .env file if exists
	godotenv.Load()

	v := viper.New()

	// Set config file path
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config")
	v.AddConfigPath(".")

	// Environment variable settings
	v.AutomaticEnv()
	v.SetEnvPrefix("SPRINGGO")

	// Default values
	setDefaults(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, use defaults and env vars
			fmt.Println("No config file found, using defaults and environment variables")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	return v, nil
}

func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.host", "localhost")

	// Database defaults
	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.name", "threadbolt_dev")
	v.SetDefault("database.username", "")
	v.SetDefault("database.password", "")
	v.SetDefault("database.sslmode", "disable")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "text")
}