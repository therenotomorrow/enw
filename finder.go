package enw

import (
	"context"
	"errors"

	"github.com/therenotomorrow/ex"
)

type (
	Source interface {
		Extract(ctx context.Context) (envs map[string]string, err error)
	}

	NamedSource struct {
		Source Source
		Name   string
	}

	Finder struct {
		storage map[string]map[string]string
		sources []NamedSource
	}
)

func NewFinder(sources []NamedSource) (*Finder, error) {
	if len(sources) == 0 {
		return nil, ErrMissingSources
	}

	uniq := make(map[string]bool)
	for _, source := range sources {
		if uniq[source.Name] {
			return nil, ErrNotUniqueSource
		}

		uniq[source.Name] = true
	}

	return &Finder{sources: sources, storage: nil}, nil
}

func (f *Finder) Find(env *Env) *Env {
	found, err := f.FindContext(context.Background(), env)

	switch {
	case errors.Is(err, ErrEnvNotFound):
	case err != nil:
		panic(ex.Cause(err))
	}

	return found
}

func (f *Finder) FindContext(ctx context.Context, env *Env) (*Env, error) {
	err := f.load(ctx)
	if err != nil {
		return nil, err
	}

	for _, source := range f.sources {
		env, ok := f.find(env, source)
		if ok {
			return env, nil
		}
	}

	return nil, ErrEnvNotFound
}

func (f *Finder) Search(env *Env) []*Env {
	return ex.Must(f.SearchContext(context.Background(), env))
}

func (f *Finder) SearchContext(ctx context.Context, env *Env) ([]*Env, error) {
	err := f.load(ctx)
	if err != nil {
		return nil, err
	}

	envs := make([]*Env, 0)

	for _, source := range f.sources {
		env, ok := f.find(env, source)
		if !ok {
			continue
		}

		envs = append(envs, env)
	}

	return envs, nil
}

func (f *Finder) find(env *Env, source NamedSource) (*Env, bool) {
	if env == nil {
		return nil, false
	}

	val, ok := f.storage[source.Name][env.Var]
	if !ok {
		return nil, false
	}

	clone := *env

	clone.Val = val
	clone.Source = source.Name

	return &clone, true
}

func (f *Finder) load(ctx context.Context) error {
	if len(f.storage) != 0 {
		return nil
	}

	f.storage = make(map[string]map[string]string)

	for _, source := range f.sources {
		data, err := source.Source.Extract(ctx)
		if err != nil {
			return ex.From(err)
		}

		f.storage[source.Name] = data
	}

	return nil
}
