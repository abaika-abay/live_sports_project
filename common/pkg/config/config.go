package config

import (
	"fmt"
	"os"
	"path/filepath" // Import for path manipulation

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl              string
	Port               string
	SportradarAPIKey   string
	SportradarBaseURL  string // New field
	NotificationBroker string
	WebSocketPort      string // New field for WebSocket server
}

func LoadConfig() (*Config, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}
	projectRoot := filepath.Dir(currentDir)
	dotEnvPath := filepath.Join(projectRoot, ".env")

	err = godotenv.Load(dotEnvPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Warning: Error loading .env file from %s: %v\n", dotEnvPath, err)
	} else if os.IsNotExist(err) {
		fmt.Println("Warning: .env file not found at project root. Attempting to load from current directory as fallback.")
		err = godotenv.Load(".env") // Try current directory as fallback
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("Warning: Error loading .env file from current directory: %v\n", err)
		} else if os.IsNotExist(err) {
			fmt.Println("Warning: No .env file found in project root or current directory. Relying on system environment variables.")
		}
	}

	cfg := &Config{
		DBUrl:              os.Getenv("DB_URL"),
		Port:               os.Getenv("PORT"),
		SportradarAPIKey:   os.Getenv("SPORTRADAR_API_KEY"),
		SportradarBaseURL:  os.Getenv("SPORTRADAR_BASE_URL"), // Get from env
		NotificationBroker: os.Getenv("NOTIFICATION_BROKER"),
		WebSocketPort:      os.Getenv("WS_PORT"), // Get WebSocket port
	}

	if cfg.DBUrl == "" {
		return nil, fmt.Errorf("DB_URL environment variable not set")
	}
	if cfg.Port == "" {
		cfg.Port = ":50051"
		fmt.Printf("PORT environment variable not set, defaulting to %s\n", cfg.Port)
	}
	if cfg.SportradarAPIKey == "" {
		fmt.Println("Warning: SPORTRADAR_API_KEY environment variable not set.")
		// For a real app, this should be an error if Sportradar is critical
	}
	if cfg.SportradarBaseURL == "" {
		cfg.SportradarBaseURL = "https://api.sportradar.us/soccer/trial/v4/en/" // Default for trial
		fmt.Printf("Warning: SPORTRADAR_BASE_URL not set, defaulting to %s\n", cfg.SportradarBaseURL)
	}
	if cfg.WebSocketPort == "" {
		cfg.WebSocketPort = ":8080" // Default WebSocket port
		fmt.Printf("Warning: WS_PORT not set, defaulting to %s\n", cfg.WebSocketPort)
	}

	return cfg, nil
}
