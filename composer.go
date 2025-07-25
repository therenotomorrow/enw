package enw

import (
	"context"
	"errors"
)

type (
	Config struct {
		Parser   Parser
		Target   any
		Sources  []NamedSource
		Autoload bool
	}
	Composer struct {
		collector *Collector
		finder    *Finder
		config    Config
	}
)

func NewComposer(config Config) (*Composer, error) {
	collector, err := NewCollector(config.Parser)
	if err != nil {
		return nil, err
	}

	finder, err := NewFinder(config.Sources)
	if err != nil {
		return nil, err
	}

	_, err = collector.Collect(config.Target)
	comp := &Composer{config: config, collector: collector, finder: finder}

	if config.Autoload {
		err = errors.Join(err, finder.load(context.Background()))
	}

	if err != nil {
		return nil, err
	}

	return comp, err
}

func (c *Composer) Collect() ([]*Env, error) {
	return c.collector.Collect(c.config.Target)
}

func (c *Composer) Find(env string) *Env {
	return c.finder.Find(New(env))
}

func (c *Composer) Search(env string) []*Env {
	return c.finder.Search(New(env))
}
