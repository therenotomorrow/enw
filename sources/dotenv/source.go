package dotenv

import (
	"context"
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/therenotomorrow/ex"
)

const (
	defaultFilename = ".env"

	ErrMissingFile ex.Const = "missing file"
)

type (
	Config struct {
		Filename string
	}

	Source struct {
		config Config
	}
)

func New() *Source {
	return NewWithConfig(Config{Filename: defaultFilename})
}

func NewWithConfig(config Config) *Source {
	if config.Filename == "" {
		config.Filename = defaultFilename
	}

	return &Source{config: config}
}

func (s *Source) Config() Config {
	return s.config
}

func (s *Source) Extract(_ context.Context) (map[string]string, error) {
	envs, err := godotenv.Read(s.config.Filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrMissingFile
	}

	if err != nil {
		return nil, ex.Unexpected(err)
	}

	return envs, nil
}
