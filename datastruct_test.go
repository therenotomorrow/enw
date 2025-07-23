package enw_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/therenotomorrow/enw"
	"github.com/therenotomorrow/ex"
)

func TestTag(t *testing.T) {
	t.Parallel()

	// `exhaustruct` + `types` testing
	_ = enw.Tag{
		Default:  "default",
		Empty:    true,
		Required: true,
	}
}

func TestEnv(t *testing.T) {
	t.Parallel()

	// `exhaustruct` + `types` testing
	_ = enw.Env{
		Field:   "field",
		Type:    "type",
		Path:    "path",
		Var:     "var",
		Val:     "val",
		Package: "package",
		Tag:     enw.Tag{Default: "default", Empty: true, Required: true},
	}
}

func TestEnvAnonymous(t *testing.T) {
	t.Parallel()

	var val enw.Env

	assert.True(t, val.Anonymous())

	val.Package = "package"

	assert.False(t, val.Anonymous())
}

func TestErrorConsistency(t *testing.T) {
	t.Parallel()

	got := make([]string, 0)
	want := []string{
		"missing target",
		"nil target",
		"invalid target, must be struct or pointer to struct",
		"missing parser",
		"missing source",
		"empty envs",
	}

	for _, err := range []ex.C{
		enw.ErrMissingTarget,
		enw.ErrNilTarget,
		enw.ErrInvalidTarget,
		enw.ErrMissingParser,
		enw.ErrMissingSource,
		enw.ErrEmptyEnvs,
	} {
		got = append(got, err.Error())
	}

	assert.Equal(t, want, got)
}
