package enw_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/therenotomorrow/enw"
)

func TestTag(t *testing.T) {
	t.Parallel()

	// `exhaustruct` + `types` testing
	var _ = enw.Tag{
		Default:  "default",
		Empty:    true,
		Required: true,
	}
}

func TestEnv(t *testing.T) {
	t.Parallel()

	// `exhaustruct` + `types` testing
	var _ = enw.Env{
		Field:   "field",
		Type:    "type",
		Path:    "path",
		Value:   "value",
		Package: "package",
		Tag:     enw.Tag{Default: "default", Empty: true, Required: true},
	}
}

func TestConstErrorError(t *testing.T) {
	t.Parallel()

	const err = enw.ConstError("sad")

	assert.Equal(t, "sad", err.Error())
}

func TestErrorConsistency(t *testing.T) {
	t.Parallel()

	got := make([]string, 0)
	want := []string{
		"missing parser",
		"missing target",
		"nil target",
		"invalid target, must be struct or pointer to struct",
	}

	for _, err := range []enw.ConstError{
		enw.ErrMissingParser,
		enw.ErrMissingTarget,
		enw.ErrNilTarget,
		enw.ErrInvalidTarget,
	} {
		got = append(got, err.Error())
	}

	assert.Equal(t, want, got)
}
