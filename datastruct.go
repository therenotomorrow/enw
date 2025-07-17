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
	Value   string
	Package string
	Tag     Tag
}

const (
	ErrMissingTarget ex.L = "missing target"
	ErrNilTarget     ex.L = "nil target"
	ErrInvalidTarget ex.L = "invalid target, must be struct or pointer to struct"
	ErrMissingParser ex.L = "missing parser"
)
