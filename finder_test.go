package enw_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/enw"
	"github.com/therenotomorrow/enw/sources/dotenv"
	"github.com/therenotomorrow/enw/sources/k8s"
	"github.com/therenotomorrow/enw/sources/memory"
	"github.com/therenotomorrow/enw/sources/system"
)

func TestNewFinder(t *testing.T) {
	t.Parallel()

	type args struct {
		sources []enw.NamedSource
	}

	tests := []struct {
		err  error
		name string
		args args
	}{
		{
			name: "success",
			args: args{sources: []enw.NamedSource{
				{Name: "system", Source: system.New()},
				{Name: "memory", Source: memory.New(nil)},
			}},
			err: nil,
		},
		{name: "failure", args: args{sources: nil}, err: enw.ErrMissingSources},
		{
			name: "not unique names",
			args: args{sources: []enw.NamedSource{
				{Name: "system", Source: system.New()},
				{Name: "system", Source: memory.New(nil)},
			}},
			err: enw.ErrNotUniqueSource,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj, err := enw.NewFinder(test.args.sources)
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

func TestSources(t *testing.T) {
	t.Parallel()

	_ = []enw.Source{
		&dotenv.Source{},
		&k8s.Source{},
		&memory.Source{},
		&system.Source{},
	}
}

func sources() []enw.NamedSource {
	return []enw.NamedSource{
		{Name: "memory1", Source: memory.New(map[string]string{"VAR_A": "val_A1", "VAR_B": "val_B1"})},
		{Name: "memory2", Source: memory.New(map[string]string{"VAR_A": "val_A2", "VAR_C": "val_C2"})},
	}
}

func TestFinderFind(t *testing.T) {
	t.Parallel()

	type args struct {
		env *enw.Env
	}

	tests := []struct {
		args args
		want *enw.Env
		name string
	}{
		{
			name: "found",
			args: args{env: enw.New("VAR_C")},
			want: &enw.Env{Var: "VAR_C", Val: "val_C2", Source: "memory2"},
		},
		{name: "not found", args: args{env: enw.New("NOT_FOUND")}, want: nil},
		{name: "nil var", args: args{env: nil}, want: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj, err := enw.NewFinder(sources())

			require.NoError(t, err)

			got := obj.Find(test.args.env)

			assert.Equal(t, test.want, got)
		})
	}

	t.Run("panics", func(t *testing.T) {
		t.Parallel()

		obj, err := enw.NewFinder(
			[]enw.NamedSource{{Name: "memory", Source: memory.New(nil).WithError(enw.ErrNilTarget)}},
		)

		require.NoError(t, err)
		assert.PanicsWithValue(t, enw.ErrNilTarget, func() { _ = obj.Find(new(enw.Env)) })
	})
}

func TestFinderFindContext(t *testing.T) {
	t.Parallel()

	type args struct {
		env *enw.Env
	}

	type want struct {
		env *enw.Env
		err error
	}

	tests := []struct {
		want want
		args args
		name string
	}{
		{
			name: "value from first source",
			args: args{env: enw.New("VAR_C")},
			want: want{env: &enw.Env{Var: "VAR_C", Val: "val_C2", Source: "memory2"}, err: nil},
		},
		{
			name: "value from second source",
			args: args{env: enw.New("VAR_C")},
			want: want{env: &enw.Env{Var: "VAR_C", Val: "val_C2", Source: "memory2"}, err: nil},
		},
		{
			name: "prioritizes first source on duplicates",
			args: args{env: enw.New("VAR_A")},
			want: want{env: &enw.Env{Var: "VAR_A", Val: "val_A1", Source: "memory1"}, err: nil},
		},
		{
			name: "var not found",
			args: args{env: enw.New("NOT_FOUND")},
			want: want{env: nil, err: enw.ErrEnvNotFound},
		},
		{
			name: "nil var",
			args: args{env: nil},
			want: want{env: nil, err: enw.ErrEnvNotFound},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj, err := enw.NewFinder(sources())

			require.NoError(t, err)

			got, err := obj.FindContext(t.Context(), test.args.env)

			require.ErrorIs(t, err, test.want.err)
			assert.Equal(t, test.want.env, got)
		})
	}

	t.Run("loading failed", func(t *testing.T) {
		t.Parallel()

		obj, err := enw.NewFinder(
			[]enw.NamedSource{{Name: "memory", Source: memory.New(nil).WithError(enw.ErrNilTarget)}},
		)

		require.NoError(t, err)

		got, err := obj.FindContext(t.Context(), new(enw.Env))

		require.ErrorIs(t, err, enw.ErrNilTarget)
		assert.Nil(t, got)
	})

	t.Run("cached loading", func(t *testing.T) {
		t.Parallel()

		obj, err := enw.NewFinder(sources())

		require.NoError(t, err)

		_, _ = obj.FindContext(t.Context(), new(enw.Env))
		_, _ = obj.FindContext(t.Context(), new(enw.Env))
	})
}

func TestFinderSearch(t *testing.T) {
	t.Parallel()

	type args struct {
		env *enw.Env
	}

	tests := []struct {
		name string
		args args
		want []*enw.Env
	}{
		{
			name: "found",
			args: args{env: enw.New("VAR_A")},
			want: []*enw.Env{
				{Var: "VAR_A", Val: "val_A1", Source: "memory1"},
				{Var: "VAR_A", Val: "val_A2", Source: "memory2"},
			},
		},
		{name: "not found", args: args{env: &enw.Env{Var: "NOT_FOUND"}}, want: []*enw.Env{}},
		{name: "nil var", args: args{env: nil}, want: []*enw.Env{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj, err := enw.NewFinder(sources())

			require.NoError(t, err)

			got := obj.Search(test.args.env)

			assert.Equal(t, test.want, got)
		})
	}

	t.Run("panics", func(t *testing.T) {
		t.Parallel()

		obj, err := enw.NewFinder(
			[]enw.NamedSource{{Name: "memory", Source: memory.New(nil).WithError(enw.ErrNilTarget)}},
		)

		require.NoError(t, err)
		assert.PanicsWithValue(t, enw.ErrNilTarget, func() { _ = obj.Search(new(enw.Env)) })
	})
}

func TestFinderSearchContext(t *testing.T) {
	t.Parallel()

	type args struct {
		env *enw.Env
	}

	type want struct {
		err  error
		envs []*enw.Env
	}

	tests := []struct {
		args args
		name string
		want want
	}{
		{
			name: "values from all sources",
			args: args{env: enw.New("VAR_A")},
			want: want{envs: []*enw.Env{
				{Var: "VAR_A", Val: "val_A1", Source: "memory1"},
				{Var: "VAR_A", Val: "val_A2", Source: "memory2"},
			}, err: nil},
		},
		{
			name: "values from first source",
			args: args{env: enw.New("VAR_B")},
			want: want{envs: []*enw.Env{
				{Var: "VAR_B", Val: "val_B1", Source: "memory1"},
			}, err: nil},
		},
		{
			name: "values from second source",
			args: args{env: enw.New("VAR_B")},
			want: want{envs: []*enw.Env{
				{Var: "VAR_B", Val: "val_B1", Source: "memory1"},
			}, err: nil},
		},
		{
			name: "vars not found",
			args: args{env: enw.New("NOT_FOUND")},
			want: want{envs: []*enw.Env{}, err: nil},
		},
		{
			name: "nil var",
			args: args{env: nil},
			want: want{envs: []*enw.Env{}, err: nil},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj, err := enw.NewFinder(sources())

			require.NoError(t, err)

			got, err := obj.SearchContext(t.Context(), test.args.env)

			require.ErrorIs(t, err, test.want.err)
			assert.Equal(t, test.want.envs, got)
		})
	}

	t.Run("loading failed", func(t *testing.T) {
		t.Parallel()

		obj, err := enw.NewFinder(
			[]enw.NamedSource{{Name: "memory", Source: memory.New(nil).WithError(enw.ErrNilTarget)}},
		)

		require.NoError(t, err)

		got, err := obj.SearchContext(t.Context(), new(enw.Env))

		require.ErrorIs(t, err, enw.ErrNilTarget)
		assert.Nil(t, got)
	})

	t.Run("cached loading", func(t *testing.T) {
		t.Parallel()

		obj, err := enw.NewFinder(sources())

		require.NoError(t, err)

		_, _ = obj.SearchContext(t.Context(), new(enw.Env))
		_, _ = obj.SearchContext(t.Context(), new(enw.Env))
	})
}
