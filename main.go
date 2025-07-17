package enw

import (
	"reflect"
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
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	var (
		currPrefix = ""
		currPath   = config.target.Type().Name()
		currPkg    = config.target.Type().PkgPath()
	)

	vars := New(config.Parser).Collect(config.target, currPrefix, currPath, currPkg)

	return vars, nil
}
