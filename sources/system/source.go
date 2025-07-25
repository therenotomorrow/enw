package system

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
	envs := make(map[string]string)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", splitN)

		if len(parts) > 1 {
			envs[parts[0]] = parts[1]
		}
	}

	return envs, nil
}
