package enw_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/enw"
	"github.com/therenotomorrow/enw/parsers/sethvargo"
	"github.com/therenotomorrow/enw/sources/memory"
	"github.com/therenotomorrow/ex"
)

func TestNewComposer(t *testing.T) {
	t.Parallel()

	type args struct {
		config enw.Config
	}

	tests := []struct {
		err  error
		name string
		args args
	}{
		{
			name: "collector error",
			args: args{config: enw.Config{}},
			err:  enw.ErrMissingParser,
		},
		{
			name: "finder error",
			args: args{config: enw.Config{Parser: sethvargo.New()}},
			err:  enw.ErrMissingSources,
		},
		{
			name: "target error without autoload",
			args: args{config: enw.Config{
				Parser:   sethvargo.New(),
				Sources:  []enw.NamedSource{{Name: "memory", Source: memory.New(nil)}},
				Autoload: false,
			}},
			err: enw.ErrNilTarget,
		},
		{
			name: "target error with autoload error (detect one)",
			args: args{config: enw.Config{
				Parser:   sethvargo.New(),
				Sources:  []enw.NamedSource{{Name: "memory", Source: memory.New(nil).WithError(enw.ErrEmptyEnvs)}},
				Autoload: true,
			}},
			err: enw.ErrNilTarget,
		},
		{
			name: "target error with autoload error (detect second)",
			args: args{config: enw.Config{
				Parser:   sethvargo.New(),
				Sources:  []enw.NamedSource{{Name: "memory", Source: memory.New(nil).WithError(enw.ErrEmptyEnvs)}},
				Autoload: true,
			}},
			err: enw.ErrEmptyEnvs,
		},
		{
			name: "autoload error",
			args: args{config: enw.Config{
				Parser:   sethvargo.New(),
				Sources:  []enw.NamedSource{{Name: "memory", Source: memory.New(nil).WithError(enw.ErrEmptyEnvs)}},
				Autoload: true,
				Target:   struct{}{},
			}},
			err: enw.ErrEmptyEnvs,
		},
		{
			name: "success",
			args: args{config: enw.Config{
				Parser:   sethvargo.New(),
				Sources:  []enw.NamedSource{{Name: "memory", Source: memory.New(nil)}},
				Autoload: true,
				Target:   struct{}{},
			}},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj, err := enw.NewComposer(test.args.config)
			if test.err != nil {
				require.ErrorIs(t, err, test.err)
				assert.Nil(t, obj)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, obj)
			}
		})
	}
}

func newComposer() *enw.Composer {
	return ex.Must(enw.NewComposer(enw.Config{
		Parser:   sethvargo.New(),
		Sources:  []enw.NamedSource{{Name: "memory", Source: memory.New(nil)}},
		Target:   struct{}{},
		Autoload: true,
	}))
}

func TestComposerCollect(t *testing.T) {
	t.Parallel()

	got, err := newComposer().Collect() // just a proxy
	want := make([]*enw.Env, 0)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestComposerFind(t *testing.T) {
	t.Parallel()

	got := newComposer().Find("mad") // just a proxy

	assert.Nil(t, got)
}

func TestComposerSearch(t *testing.T) {
	t.Parallel()

	got := newComposer().Search("mad") // just a proxy
	want := make([]*enw.Env, 0)

	assert.Equal(t, want, got)
}
