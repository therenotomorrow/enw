package system_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/enw/sources/system"
)

func TestNew(t *testing.T) {
	t.Parallel()

	obj := system.New()

	assert.NotNil(t, obj)
}

func TestSourceExtract(t *testing.T) {
	t.Setenv("TestSourceExtract", "no-parallel")

	obj := system.New()

	tests := []struct {
		setup func()
		want  map[string]string
		name  string
	}{
		{
			name: "standard environment variables",
			setup: func() {
				t.Setenv("KEY1", "VALUE1")
				t.Setenv("KEY2", "VALUE2")
			},
			want: map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2",
			},
		},
		{
			name: "variable with empty value",
			setup: func() {
				t.Setenv("EMPTY_VAL", "")
			},
			want: map[string]string{
				"EMPTY_VAL": "",
			},
		},
		{
			name: "variable value contains an equals sign",
			setup: func() {
				t.Setenv("COMPLEX_VAL", "value=with=equals")
			},
			want: map[string]string{
				"COMPLEX_VAL": "value=with=equals",
			},
		},
		{
			name:  "no environment variables set",
			setup: func() {},
			want:  map[string]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("TestSourceExtract."+test.name, "no-parallel")

			test.setup()

			got, err := obj.Extract(t.Context())

			fot := make(map[string]string)

			for k := range test.want {
				if val, ok := got[k]; ok {
					fot[k] = val
				}
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, fot)
		})
	}
}
