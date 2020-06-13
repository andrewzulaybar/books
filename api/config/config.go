package config

import (
	"path"
	"runtime"

	"github.com/joho/godotenv"
)

// Config holds the configuration variables needed to start the application.
type Config struct {
	ConnectionString string
	Address          string
}

// Load returns the environment variables set in the given file.
func Load(fileName string) (*Config, error) {
	_, file, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(file), fileName)

	env, err := godotenv.Read(path)
	if err != nil {
		return nil, err
	}

	return &Config{
		ConnectionString: env["CONNECTION_STRING"],
		Address:          env["ADDRESS"],
	}, nil
}
