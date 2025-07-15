package enw

import (
	"cmp"
	"reflect"
	"slices"
)

type Config struct {
	Target any
	Parser Parser
	// for internal usage
	target reflect.Value
}

func (c *Config) Validate() error {
	if c.Target == nil {
		return ErrMissingTarget
	}

	if c.Parser == nil {
		return ErrMissingParser
	}

	c.target = reflect.ValueOf(c.Target)

	if c.target.Kind() == reflect.Ptr {
		if c.target.IsNil() {
			return ErrNilTarget
		}

		c.target = c.target.Elem()
	}

	if c.target.Kind() != reflect.Struct {
		return ErrInvalidTarget
	}

	return nil
}

func Collect(config Config) ([]*Env, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	var (
		currPrefix = ""
		currPath   = config.target.Type().Name()
		currPkg    = config.target.Type().PkgPath()
	)

	vars := New(config.Parser).Collect(config.target, currPrefix, currPath, currPkg)

	slices.SortStableFunc(vars, func(a, b *Env) int {
		return cmp.Compare(a.Value, b.Value)
	})

	return vars, nil
}
