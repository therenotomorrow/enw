package dotenv

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/therenotomorrow/ex"
)

const (
	defaultFilename = ".env"

	ErrMissingFile ex.C = "missing file"
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
	mapping, err := godotenv.Read(s.config.Filename)
	if err != nil {
		return nil, ErrMissingFile.Because(err)
	}

	return mapping, nil
}
