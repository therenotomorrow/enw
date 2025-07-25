package memory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/enw/sources/memory"
	"github.com/therenotomorrow/ex"
)

func TestNew(t *testing.T) {
	t.Parallel()

	obj1 := memory.New(nil)
	obj2 := memory.New(map[string]string{"test": "test", "same": "same"})

	assert.NotNil(t, obj1)
	assert.NotNil(t, obj2)
}

func TestSourceExtract(t *testing.T) {
	t.Parallel()

	const dummyErr = ex.Const("dummy error")

	data := map[string]string{"test": "test", "same": "same"}

	t.Run("with nil data", func(t *testing.T) {
		t.Parallel()

		got, err := memory.New(nil).Extract(t.Context())

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("with nil data and error", func(t *testing.T) {
		t.Parallel()

		got, err := memory.New(nil).WithError(dummyErr).Extract(t.Context())

		require.ErrorIs(t, err, dummyErr)
		assert.Empty(t, got)
	})

	t.Run("with real data", func(t *testing.T) {
		t.Parallel()

		got, err := memory.New(data).Extract(t.Context())

		require.NoError(t, err)
		assert.Equal(t, data, got)
	})

	t.Run("with real data and error", func(t *testing.T) {
		t.Parallel()

		got, err := memory.New(data).WithError(dummyErr).Extract(t.Context())

		require.ErrorIs(t, err, dummyErr)
		assert.Equal(t, data, got)
	})
}

func TestSourceWithError(t *testing.T) {
	t.Parallel()

	obj1 := memory.New(nil)
	obj2 := obj1.WithError(nil)

	assert.NotSame(t, obj1, obj2)
}
