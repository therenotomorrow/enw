package dotenv

import (
	"context"
	"os"
	"strings"
)

const (
	splitN = 2
)

type Source struct{}

func New() *Source {
	return &Source{}
}

func (s *Source) Extract(_ context.Context) (map[string]string, error) {
	mapping := make(map[string]string)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", splitN)
		mapping[parts[0]] = parts[1]
	}

	return mapping, nil
}
