package sethvargo_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/therenotomorrow/enw"
	"github.com/therenotomorrow/enw/parsers/sethvargo"
)

func TestNew(t *testing.T) {
	t.Parallel()

	obj := sethvargo.New()

	assert.Equal(t, sethvargo.Config{TagKey: "env"}, obj.Config())
}

func TestNewWithConfig(t *testing.T) {
	t.Parallel()

	type args struct {
		config sethvargo.Config
	}

	tests := []struct {
		name string
		args args
		want sethvargo.Config
	}{
		{
			name: "default key",
			args: args{config: sethvargo.Config{}},
			want: sethvargo.Config{TagKey: "env"},
		},
		{
			name: "custom key",
			args: args{config: sethvargo.Config{TagKey: "custom"}},
			want: sethvargo.Config{TagKey: "custom"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj := sethvargo.NewWithConfig(test.args.config)

			assert.Equal(t, test.want, obj.Config())
		})
	}
}

func TestParserParse(t *testing.T) {
	t.Parallel()

	type sampleStruct struct {
		Simple           string    `env:"MY_VAR"`
		WithDefault      string    `env:"MY_VAR,default=fallback"`
		WithRequired     string    `env:"MY_VAR,required"`
		WithPrefix       string    `env:"MY_VAR,prefix=APP_"`
		WithAllOptions   string    `env:"MY_VAR,default=fallback,required,prefix=APP_"`
		WithSpaces       string    `env:"  MY_VAR  ,  default=fallback ,required, prefix=APP_ "`
		OnlyPrefix       string    `env:",prefix=APP_"`
		OnlyDefault      string    `env:",default=fallback"`
		OnlyRequired     string    `env:",required"`
		EmptyTag         string    `env:""`
		IgnoredTag       string    `env:"-"`
		NoEnvTag         string    `custom:"tag"`
		OnlyVariableName string    `env:"MY_VAR"`
		EmptyDefault     string    `env:"MY_VAR,default="`
		JustAComma       string    `env:","`
		ExternalType     time.Time `env:"TIME_VAR"`
	}

	type want struct {
		env    *enw.Env
		prefix string
	}

	parser := sethvargo.New()
	sample := sampleStruct{}

	tests := []struct {
		name  string
		field string
		want  want
	}{
		{
			name:  "simple var",
			field: "Simple",
			want: want{
				env: &enw.Env{
					Value:   "MY_VAR",
					Field:   "Simple",
					Type:    "string",
					Path:    "some.path",
					Package: "some/pkg",
					Tag:     enw.Tag{Default: "", Required: false, Empty: true},
				},
				prefix: "",
			},
		},
		{
			name:  "with default",
			field: "WithDefault",
			want: want{
				env: &enw.Env{
					Value:   "MY_VAR",
					Field:   "WithDefault",
					Type:    "string",
					Path:    "some.path",
					Package: "some/pkg",
					Tag:     enw.Tag{Default: "fallback", Required: false, Empty: false},
				},
				prefix: "",
			},
		},
		{
			name:  "with required",
			field: "WithRequired",
			want: want{
				env: &enw.Env{
					Value:   "MY_VAR",
					Field:   "WithRequired",
					Type:    "string",
					Path:    "some.path",
					Package: "some/pkg",
					Tag:     enw.Tag{Default: "", Required: true, Empty: false},
				},
				prefix: "",
			},
		},
		{
			name:  "with prefix",
			field: "WithPrefix",
			want: want{
				env: &enw.Env{
					Value:   "MY_VAR",
					Field:   "WithPrefix",
					Type:    "string",
					Path:    "some.path",
					Package: "some/pkg",
					Tag:     enw.Tag{Default: "", Required: false, Empty: false},
				},
				prefix: "APP_",
			},
		},
		{
			name:  "with all options",
			field: "WithAllOptions",
			want: want{
				env: &enw.Env{
					Value:   "MY_VAR",
					Field:   "WithAllOptions",
					Type:    "string",
					Path:    "some.path",
					Package: "some/pkg",
					Tag:     enw.Tag{Default: "fallback", Required: true, Empty: false},
				},
				prefix: "APP_",
			},
		},
		{
			name:  "with spaces",
			field: "WithSpaces",
			want: want{
				env: &enw.Env{
					Value:   "MY_VAR",
					Field:   "WithSpaces",
					Type:    "string",
					Path:    "some.path",
					Package: "some/pkg",
					Tag:     enw.Tag{Default: "fallback", Required: true, Empty: false},
				},
				prefix: "APP_",
			},
		},
		{
			name:  "empty default value",
			field: "EmptyDefault",
			want: want{
				env: &enw.Env{
					Value:   "MY_VAR",
					Field:   "EmptyDefault",
					Type:    "string",
					Path:    "some.path",
					Package: "some/pkg",
					Tag:     enw.Tag{Default: "", Required: false, Empty: true},
				},
				prefix: "",
			},
		},
		{
			name:  "field with type from external package",
			field: "ExternalType",
			want: want{
				env: &enw.Env{
					Value:   "TIME_VAR",
					Field:   "ExternalType",
					Type:    "time.Time",
					Path:    "some.path",
					Package: "some/pkg",
					Tag:     enw.Tag{Default: "", Required: false, Empty: true},
				},
				prefix: "",
			},
		},
		{name: "only prefix", field: "OnlyPrefix", want: want{env: nil, prefix: "APP_"}},
		{name: "only default", field: "OnlyDefault", want: want{env: nil, prefix: ""}},
		{name: "only required", field: "OnlyRequired", want: want{env: nil, prefix: ""}},
		{name: "empty tag", field: "EmptyTag", want: want{env: nil, prefix: ""}},
		{name: "ignored tag", field: "IgnoredTag", want: want{env: nil, prefix: ""}},
		{name: "no env tag", field: "NoEnvTag", want: want{env: nil, prefix: ""}},
		{name: "just a comma", field: "JustAComma", want: want{env: nil, prefix: ""}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			structField, ok := reflect.TypeOf(sample).FieldByName(test.field)

			assert.True(t, ok)

			got, prefix := parser.Parse(structField, "some.path", "some/pkg")

			assert.Equal(t, test.want.env, got)
			assert.Equal(t, test.want.prefix, prefix)
		})
	}
}
