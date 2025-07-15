package enw

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

type ConstError string

func (e ConstError) Error() string {
	return string(e)
}

const (
	ErrMissingTarget ConstError = "missing target"
	ErrNilTarget     ConstError = "nil target"
	ErrInvalidTarget ConstError = "invalid target, must be struct or pointer to struct"
	ErrMissingParser ConstError = "missing parser"
)
