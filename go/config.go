package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
	PokeAPI  PokeAPIConfig  `mapstructure:"pokeapi"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type PokeAPIConfig struct {
	URL string `mapstructure:"url"`
	Max int    `mapstructure:"max"`
}

// LoadConfig reads configuration from environment variables and config file
func LoadConfig() (*Config, error) {
	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("pokeapi.url", "https://pokeapi.co/api/v2/pokemon?limit=")
	viper.SetDefault("pokeapi.max", 1025)

	// Environment variable support
	viper.SetEnvPrefix("POKE")
	viper.AutomaticEnv()

	// Try to read config file
	viper.SetConfigFile(".env")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		// Config file is optional, only log if it exists but can't be read
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Override with environment variables if they exist
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		viper.Set("database.url", dbURL)
	}
	if port := os.Getenv("PORT"); port != "" {
		viper.Set("server.port", port)
	}
	if pokeAPIURL := os.Getenv("POKEAPI_URL"); pokeAPIURL != "" {
		viper.Set("pokeapi.url", pokeAPIURL)
	}
	if pokeAPIMaxStr := os.Getenv("POKEAPI_MAX"); pokeAPIMaxStr != "" {
		if max, err := strconv.Atoi(pokeAPIMaxStr); err == nil {
			viper.Set("pokeapi.max", max)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.PokeAPI.URL == "" {
		return fmt.Errorf("PokeAPI URL is required")
	}
	if c.PokeAPI.Max <= 0 {
		return fmt.Errorf("PokeAPI max must be greater than 0")
	}
	return nil
}