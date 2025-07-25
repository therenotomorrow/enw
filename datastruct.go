package enw

import (
	"github.com/therenotomorrow/ex"
)

type Tag struct {
	Default  string
	Empty    bool
	Required bool
}

type Env struct {
	Field   string
	Type    string
	Path    string
	Var     string
	Val     string
	Package string
	Source  string
	Tag     Tag
}

func New(key string) *Env {
	env := new(Env)
	env.Var = key

	return env
}

func (e *Env) Anonymous() bool {
	return e.Package == ""
}

const (
	ErrMissingTarget   ex.Const = "missing target"
	ErrNilTarget       ex.Const = "nil target"
	ErrInvalidTarget   ex.Const = "invalid target, must be struct or pointer to struct"
	ErrMissingParser   ex.Const = "missing parser"
	ErrMissingSources  ex.Const = "missing sources"
	ErrNotUniqueSource ex.Const = "not unique source"
	ErrEmptyEnvs       ex.Const = "empty envs"
	ErrEnvNotFound     ex.Const = "env not found"
)
