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
	Tag     Tag
}

func (e *Env) Anonymous() bool {
	return e.Package == ""
}

const (
	ErrMissingTarget ex.C = "missing target"
	ErrNilTarget     ex.C = "nil target"
	ErrInvalidTarget ex.C = "invalid target, must be struct or pointer to struct"
	ErrMissingParser ex.C = "missing parser"
	ErrMissingSource ex.C = "missing source"
	ErrEmptyEnvs     ex.C = "empty envs"
)
