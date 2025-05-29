package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	Mongo struct {
		URI      string
		Database string
	}
	Log struct {
		Level string
	}
	NATS struct {
		URL string
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../common/pkg/config") // For flexibility

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Override with environment variables if set
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		cfg.Mongo.URI = uri
	}
	// Add more overrides as needed

	return &cfg, nil
}
