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
		Source:  "source",
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
		"missing sources",
		"empty envs",
		"not unique source",
	}

	for _, err := range []ex.Const{
		enw.ErrMissingTarget,
		enw.ErrNilTarget,
		enw.ErrInvalidTarget,
		enw.ErrMissingParser,
		enw.ErrMissingSources,
		enw.ErrEmptyEnvs,
		enw.ErrNotUniqueSource,
	} {
		got = append(got, err.Error())
	}

	assert.Equal(t, want, got)
}

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "empty", args: args{key: ""}, want: ""},
		{name: "smoke", args: args{key: "hehe"}, want: "hehe"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := enw.New(test.args.key)
			obj := new(enw.Env)

			obj.Var = test.args.key

			assert.Equal(t, test.want, got.Var)
			assert.Equal(t, obj, got)
		})
	}
}
