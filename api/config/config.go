package config

import (
	"github.com/joho/godotenv"
)

// Config holds the configuration variables needed to start the application.
type Config struct {
	ConnectionString string
	Port             string
}

// Load returns the environment variables set in the file of the given path.
func Load(path string) (*Config, error) {
	env, err := godotenv.Read(path)
	if err != nil {
		return nil, err
	}

	return &Config{
		ConnectionString: env["CONNECTION_STRING"],
		Port:             env["PORT"],
	}, nil
}
