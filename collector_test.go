package enw_test

import (
	"cmp"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/enw"
	"github.com/therenotomorrow/enw/parsers/sethvargo"
)

const (
	testPackage = "github.com/therenotomorrow/enw_test"
)

func TestNewCollector(t *testing.T) {
	t.Parallel()

	type args struct {
		parser enw.Parser
	}

	tests := []struct {
		args args
		err  error
		name string
	}{
		{name: "success", args: args{parser: sethvargo.New()}, err: nil},
		{name: "failure", args: args{parser: nil}, err: enw.ErrMissingParser},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj, err := enw.NewCollector(test.args.parser)
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

func TestParsers(t *testing.T) {
	t.Parallel()

	_ = []enw.Parser{
		&sethvargo.Parser{},
	}
}

func wantCollectorCollect() []*enw.Env {
	return []*enw.Env{
		{
			Var:     "APP_NAME",
			Field:   "AppName",
			Type:    "string",
			Path:    "sampleConfig->AppName",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "DB_HOST",
			Field:   "Host",
			Type:    "string",
			Path:    "sampleConfig->DB->Host",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "DB_PORT",
			Field:   "Port",
			Type:    "int",
			Path:    "sampleConfig->DB->Port",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "CACHE_HOST",
			Field:   "Host",
			Type:    "string",
			Path:    "sampleConfig->Cache->Host",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "CACHE_PORT",
			Field:   "Port",
			Type:    "int",
			Path:    "sampleConfig->Cache->Port",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "SRV_HOST",
			Field:   "Host",
			Type:    "string",
			Path:    "sampleConfig->Servers->0->Host",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "SRV_PORT",
			Field:   "Port",
			Type:    "int",
			Path:    "sampleConfig->Servers->0->Port",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "SRV_HOST",
			Field:   "Host",
			Type:    "string",
			Path:    "sampleConfig->Servers->1->Host",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "SRV_PORT",
			Field:   "Port",
			Type:    "int",
			Path:    "sampleConfig->Servers->1->Port",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "PTR_SRV_HOST",
			Field:   "Host",
			Type:    "string",
			Path:    "sampleConfig->PtrServers->0->Host",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "PTR_SRV_PORT",
			Field:   "Port",
			Type:    "int",
			Path:    "sampleConfig->PtrServers->0->Port",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "PTR_SRV_HOST",
			Field:   "Host",
			Type:    "string",
			Path:    "sampleConfig->PtrServers->1->Host",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "PTR_SRV_PORT",
			Field:   "Port",
			Type:    "int",
			Path:    "sampleConfig->PtrServers->1->Port",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "NIL_SRV_HOST",
			Field:   "Host",
			Type:    "string",
			Path:    "sampleConfig->NilInSlice->0->Host",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "NIL_SRV_PORT",
			Field:   "Port",
			Type:    "int",
			Path:    "sampleConfig->NilInSlice->0->Port",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "NIL_SRV_HOST",
			Field:   "Host",
			Type:    "string",
			Path:    "sampleConfig->NilInSlice->2->Host",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
		{
			Var:     "NIL_SRV_PORT",
			Field:   "Port",
			Type:    "int",
			Path:    "sampleConfig->NilInSlice->2->Port",
			Package: testPackage,
			Tag:     enw.Tag{Empty: true},
		},
	}
}

func TestCollectorCollect(t *testing.T) {
	t.Parallel()

	type args struct {
		target any
	}

	type want struct {
		err  error
		envs []*enw.Env
	}

	type Sample struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
	}

	type sampleConfig struct {
		Cache      *Sample   `env:",prefix=CACHE_"`
		EmptyCache *Sample   `env:",prefix=EMPTY_"`
		AppName    string    `env:"APP_NAME"`
		unexported string    `env:"UNEXPORTED"` //nolint:unused // this is used for tests
		DB         Sample    `env:",prefix=DB_"`
		Servers    []Sample  `env:",prefix=SRV_"`
		PtrServers []*Sample `env:",prefix=PTR_SRV_"`
		NilInSlice []*Sample `env:",prefix=NIL_SRV_"`
	}

	var (
		srv1      = Sample{Host: "srv1.local", Port: 8080}
		srv2      = Sample{Host: "srv2.local", Port: 8081}
		dbConf    = Sample{Host: "db.local", Port: 5432}
		cacheConf = Sample{Host: "cache.local", Port: 6379}
	)

	tests := []struct {
		args args
		name string
		want want
	}{
		{
			name: "full config walk as a struct",
			args: args{target: sampleConfig{
				AppName:    "MyApp",
				DB:         dbConf,
				Cache:      &cacheConf,
				EmptyCache: nil,
				Servers:    []Sample{srv1, srv2},
				PtrServers: []*Sample{&srv1, &srv2},
				NilInSlice: []*Sample{&srv1, nil, &srv2},
			}},
			want: want{envs: wantCollectorCollect(), err: nil},
		},
		{
			name: "full config walk as a pointer",
			args: args{target: &sampleConfig{
				AppName:    "MyApp",
				DB:         dbConf,
				Cache:      &cacheConf,
				EmptyCache: nil,
				Servers:    []Sample{srv1, srv2},
				PtrServers: []*Sample{&srv1, &srv2},
				NilInSlice: []*Sample{&srv1, nil, &srv2},
			}},
			want: want{envs: wantCollectorCollect(), err: nil},
		},
		{
			name: "nil struct",
			args: args{target: nil},
			want: want{envs: nil, err: enw.ErrNilTarget},
		},
		{
			name: "typed but nil struct",
			args: args{target: (*sampleConfig)(nil)},
			want: want{envs: nil, err: enw.ErrNilTarget},
		},
		{
			name: "empty struct",
			args: args{target: struct{}{}},
			want: want{envs: []*enw.Env{}, err: nil},
		},
		{
			name: "not a struct",
			args: args{target: 123},
			want: want{envs: nil, err: enw.ErrInvalidTarget},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj, err := enw.NewCollector(sethvargo.New())

			require.NoError(t, err)

			got, err := obj.Collect(test.args.target)

			require.ErrorIs(t, err, test.want.err)

			slices.SortStableFunc(test.want.envs, func(a, b *enw.Env) int {
				return cmp.Compare(a.Var, b.Var)
			})

			assert.Len(t, got, len(test.want.envs))
			assert.Equal(t, test.want.envs, got)
		})
	}
}
