package enw

import (
	"context"
	"sync"

	"github.com/therenotomorrow/ex"
)

type (
	Source interface {
		Extract(ctx context.Context) (varToVal map[string]string, err error)
	}

	Finder struct {
		source  Source
		mapping map[string]string
		mutex   sync.Mutex
	}
)

func NewFinder(source Source) (*Finder, error) {
	if source == nil {
		return nil, ErrMissingSource
	}

	return &Finder{mutex: sync.Mutex{}, source: source, mapping: nil}, nil
}

func (f *Finder) Find(envs []*Env) (map[string]string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(envs) == 0 {
		return nil, ErrEmptyEnvs
	}

	if f.mapping != nil {
		return f.mapping, nil
	}

	source, err := f.source.Extract(context.Background())
	if err != nil {
		return nil, ex.From(err)
	}

	f.mapping = make(map[string]string)

	for _, env := range envs {
		f.mapping[env.Var] = source[env.Var]
	}

	return f.mapping, nil
}
