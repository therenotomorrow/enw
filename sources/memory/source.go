package memory

import (
	"context"
)

type Source struct {
	data map[string]string
	err  error
}

func (s *Source) WithError(err error) *Source {
	return &Source{data: s.data, err: err}
}

func New(data map[string]string) *Source {
	if data == nil {
		data = make(map[string]string)
	}

	return &Source{data: data, err: nil}
}

func (s *Source) Extract(_ context.Context) (map[string]string, error) {
	return s.data, s.err
}
