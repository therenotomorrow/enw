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

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	type sampleConfig struct {
		Field string `env:"FIELD"`
	}

	tests := []struct {
		want   error
		config enw.Config
		name   string
	}{
		{
			name:   "success with struct",
			config: enw.Config{Target: sampleConfig{}, Parser: sethvargo.New()},
			want:   nil,
		},
		{
			name:   "success with pointer",
			config: enw.Config{Target: &sampleConfig{}, Parser: sethvargo.New()},
			want:   nil,
		},
		{
			name:   "error on nil target",
			config: enw.Config{Target: nil, Parser: sethvargo.New()},
			want:   enw.ErrMissingTarget,
		},
		{
			name:   "error on nil parser",
			config: enw.Config{Target: sampleConfig{}, Parser: nil},
			want:   enw.ErrMissingParser,
		},
		{
			name:   "error on nil target pointer",
			config: enw.Config{Target: (*sampleConfig)(nil), Parser: sethvargo.New()},
			want:   enw.ErrNilTarget,
		},
		{
			name:   "error on invalid target (int)",
			config: enw.Config{Target: 123, Parser: sethvargo.New()},
			want:   enw.ErrInvalidTarget,
		},
		{
			name:   "error on invalid target pointer (int)",
			config: enw.Config{Target: new(int), Parser: sethvargo.New()},
			want:   enw.ErrInvalidTarget,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.config.Validate()

			require.ErrorIs(t, err, test.want)
		})
	}
}

func wantCollect() []*enw.Env {
	return []*enw.Env{
		{
			Value:   "APP_NAME",
			Field:   "AppName",
			Type:    "string",
			Path:    "sampleConfig->AppName",
			Package: testPackage,
			Tag:     enw.Tag{Default: "", Empty: true, Required: false},
		},
		{
			Value:   "DB_HOST",
			Field:   "Host",
			Type:    "string",
			Path:    "sampleConfig->DB->Host",
			Package: testPackage,
			Tag:     enw.Tag{Default: "", Empty: true, Required: false},
		},
		{
			Value:   "DB_PORT",
			Field:   "Port",
			Type:    "int",
			Path:    "sampleConfig->DB->Port",
			Package: testPackage,
			Tag:     enw.Tag{Default: "", Empty: true, Required: false},
		},
		{
			Value:   "SRV_HOST",
			Field:   "Host",
			Type:    "string",
			Path:    "sampleConfig->Servers->0->Host",
			Package: testPackage,
			Tag:     enw.Tag{Default: "", Empty: true, Required: false},
		},
		{
			Value:   "SRV_PORT",
			Field:   "Port",
			Type:    "int",
			Path:    "sampleConfig->Servers->0->Port",
			Package: testPackage,
			Tag:     enw.Tag{Default: "", Empty: true, Required: false},
		},
	}
}

func TestCollect(t *testing.T) {
	t.Parallel()

	type Sample struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
	}

	type sampleConfig struct {
		AppName string   `env:"APP_NAME"`
		DB      Sample   `env:",prefix=DB_"`
		Servers []Sample `env:",prefix=SRV_"`
	}

	type want struct {
		err  error
		envs []*enw.Env
	}

	var (
		srv1   = Sample{Host: "srv1.local", Port: 8080}
		dbConf = Sample{Host: "db.local", Port: 5432}
	)

	tests := []struct {
		config enw.Config
		name   string
		want   want
	}{
		{
			name: "Success with Struct Value",
			config: enw.Config{
				Target: sampleConfig{AppName: "MyApp", DB: dbConf, Servers: []Sample{srv1}},
				Parser: sethvargo.New(),
			},
			want: want{envs: wantCollect(), err: nil},
		},
		{
			name: "Success with Struct Pointer",
			config: enw.Config{
				Target: &sampleConfig{AppName: "MyApp", DB: dbConf, Servers: []Sample{srv1}},
				Parser: sethvargo.New(),
			},
			want: want{envs: wantCollect(), err: nil},
		},
		{
			name: "Failure on Invalid Target",
			config: enw.Config{
				Target: 123,
				Parser: sethvargo.New(),
			},
			want: want{envs: nil, err: enw.ErrInvalidTarget},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := enw.Collect(test.config)

			slices.SortStableFunc(test.want.envs, func(a, b *enw.Env) int {
				return cmp.Compare(a.Value, b.Value)
			})

			require.ErrorIs(t, err, test.want.err)
			assert.Equal(t, test.want.envs, got)
		})
	}
}
