package dotenv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/enw/sources/dotenv"
	"github.com/therenotomorrow/ex"
)

func TestNew(t *testing.T) {
	t.Parallel()

	obj := dotenv.New()

	assert.Equal(t, dotenv.Config{Filename: ".env"}, obj.Config())
}

func TestNewWithConfig(t *testing.T) {
	t.Parallel()

	type args struct {
		config dotenv.Config
	}

	tests := []struct {
		name string
		args args
		want dotenv.Config
	}{
		{
			name: "default filename",
			args: args{config: dotenv.Config{}},
			want: dotenv.Config{Filename: ".env"},
		},
		{
			name: "custom filename",
			args: args{config: dotenv.Config{Filename: ".env.example"}},
			want: dotenv.Config{Filename: ".env.example"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj := dotenv.NewWithConfig(test.args.config)

			assert.Equal(t, test.want, obj.Config())
		})
	}
}

func testFile(t *testing.T, content string) string {
	t.Helper()

	file := ex.Must(os.CreateTemp(t.TempDir(), ".env.*.test"))

	_ = ex.Must(file.WriteString(content))
	ex.MustDo(file.Close())

	return file.Name()
}

func TestSourceExtract(t *testing.T) {
	t.Parallel()

	type args struct {
		content string
	}

	type want struct {
		envs map[string]string
		err  error
	}

	tests := []struct {
		want want
		name string
		args args
	}{
		{
			name: "success",
			args: args{content: "KEY1=VALUE1\nKEY2=VALUE2"},
			want: want{envs: map[string]string{"KEY1": "VALUE1", "KEY2": "VALUE2"}, err: nil},
		},
		{
			name: "some broken content",
			args: args{content: "KEY1\nKEY2===VALUE2\nKEY3={\n}"},
			want: want{envs: nil, err: ex.ErrUnexpected},
		},
		{
			name: "file not found",
			args: args{content: ""},
			want: want{envs: nil, err: dotenv.ErrMissingFile},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var filename string
			if test.args.content != "" {
				filename = testFile(t, test.args.content)
			}

			obj := dotenv.NewWithConfig(dotenv.Config{Filename: filename})

			got, err := obj.Extract(t.Context())

			require.ErrorIs(t, err, test.want.err)
			assert.Equal(t, test.want.envs, got)
		})
	}
}
