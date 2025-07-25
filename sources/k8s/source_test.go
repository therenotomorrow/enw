package k8s_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/enw/sources/k8s"
	"github.com/therenotomorrow/ex"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err    error
		config k8s.Config
		name   string
	}{
		{name: "valid configmap", config: k8s.Config{Name: "configmap", Type: k8s.ConfigMap}, err: nil},
		{name: "valid secret", config: k8s.Config{Name: "secret", Type: k8s.Secret}, err: nil},
		{name: "missing name", config: k8s.Config{Type: k8s.ConfigMap}, err: k8s.ErrMissingName},
		{name: "missing type", config: k8s.Config{Name: "name"}, err: k8s.ErrMissingType},
		{name: "invalid type", config: k8s.Config{Name: "name", Type: "invalid"}, err: k8s.ErrInvalidType},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.config.Validate()

			require.ErrorIs(t, err, test.err)
		})
	}
}

const (
	failureKonfig = `
apiVersion: v1
kind: Config
clusters:
  - cluster:
      server: https://any-server.local
    name: cluster
contexts:
  - context:
      cluster: cluster
      user: user
    name: context
current-context: context
preferences: {}
users:
  - name: user
    user:
      client-certificate-data: QkFTRTY0RU5DT0RFRF9DTElFTlRfQ0VSVAo=
      client-key-data: QkFTRTY0RU5DT0RFRF9DTElFTlRfQ0VSVAo=`

	successKonfig = `
apiVersion: v1
kind: Config
clusters:
  - cluster:
      server: http://127.0.0.1
    name: cluster-1
  - cluster:
      server: http://127.0.0.2
    name: cluster-2
contexts:
  - context:
      cluster: cluster-1
      namespace: myspace
      user: user-1
    name: context-1
  - context:
      cluster: cluster-2
      user: user-2
    name: context-2
  - context:
      cluster: cluster-1
      namespace: spacer
      user: user-2
    name: context-3
current-context: context-1
preferences: {}
users:
  - name: user-1
    user: {}
  - name: user-2
    user: {}`

	context1 = "context-1"
	context2 = "context-2"
	context3 = "context-3"
)

func testFile(t *testing.T, content string) string {
	t.Helper()

	file := ex.Must(os.CreateTemp(t.TempDir(), "kube.*.yaml"))

	_ = ex.Must(file.WriteString(content))
	ex.MustDo(file.Close())

	return file.Name()
}

func TestNewWithConfigFailure(t *testing.T) {
	t.Setenv("TestNewWithConfigFailure", "no-parallel")

	cfg := k8s.Config{Name: "name", Namespace: "", Type: k8s.Secret, Context: ""}

	t.Run("invalid config", func(t *testing.T) {
		t.Setenv("KUBECONFIG", t.TempDir())

		var missing k8s.Config

		obj, err := k8s.NewWithConfig(missing)

		require.ErrorIs(t, err, k8s.ErrMissingName)
		assert.Nil(t, obj)
	})

	t.Run("load raw config error", func(t *testing.T) {
		t.Setenv("KUBECONFIG", t.TempDir())

		obj, err := k8s.NewWithConfig(cfg)

		require.ErrorIs(t, err, k8s.ErrKubectlError)
		require.ErrorContains(t, err, "error loading config file")
		assert.Nil(t, obj)
	})

	t.Run("load client config error", func(t *testing.T) {
		t.Setenv("KUBECONFIG", "not-exist.yaml")

		obj, err := k8s.NewWithConfig(cfg)

		require.ErrorIs(t, err, k8s.ErrKubectlError)
		require.ErrorContains(t, err, "no configuration has been provided")
		assert.Nil(t, obj)
	})

	t.Run("load client config error", func(t *testing.T) {
		t.Setenv("KUBECONFIG", "not-exist.yaml")

		obj, err := k8s.NewWithConfig(cfg)

		require.ErrorIs(t, err, k8s.ErrKubectlError)
		require.ErrorContains(t, err, "no configuration has been provided")
		assert.Nil(t, obj)
	})

	t.Run("load client error", func(t *testing.T) {
		t.Setenv("KUBECONFIG", testFile(t, failureKonfig))

		obj, err := k8s.NewWithConfig(cfg)

		require.ErrorIs(t, err, k8s.ErrKubectlError)
		require.ErrorContains(t, err, "failed to find any PEM data in certificate input")
		assert.Nil(t, obj)
	})
}

func TestNewWithConfigSuccess(t *testing.T) {
	t.Setenv("TestNewWithConfigSuccess", "no-parallel")
	t.Setenv("KUBECONFIG", testFile(t, successKonfig))

	t.Run("load without context", func(t *testing.T) {
		t.Setenv("TestNewWithConfigSuccess."+t.Name(), "no-parallel")

		cfg := k8s.Config{Name: "name", Namespace: "", Type: k8s.Secret, Context: ""}

		obj, err := k8s.NewWithConfig(cfg)

		cfg.Namespace = "myspace"
		cfg.Context = context1

		require.NoError(t, err)
		assert.Equal(t, cfg, obj.Config())
	})

	t.Run("load without context but with namespace", func(t *testing.T) {
		t.Setenv("TestNewWithConfigSuccess."+t.Name(), "no-parallel")

		cfg := k8s.Config{Name: "name", Namespace: "my", Type: k8s.Secret, Context: ""}

		obj, err := k8s.NewWithConfig(cfg)

		cfg.Namespace = "my"
		cfg.Context = context1

		require.NoError(t, err)
		assert.Equal(t, cfg, obj.Config())
	})

	t.Run("load with own context", func(t *testing.T) {
		t.Setenv("TestNewWithConfigSuccess."+t.Name(), "no-parallel")

		cfg := k8s.Config{Name: "name", Namespace: "", Type: k8s.Secret, Context: context3}

		obj, err := k8s.NewWithConfig(cfg)

		cfg.Namespace = "spacer"
		cfg.Context = context3

		require.NoError(t, err)
		assert.Equal(t, cfg, obj.Config())
	})

	t.Run("load with own context without namespace", func(t *testing.T) {
		t.Setenv("TestNewWithConfigSuccess."+t.Name(), "no-parallel")

		cfg := k8s.Config{Name: "name", Namespace: "", Type: k8s.Secret, Context: context2}

		obj, err := k8s.NewWithConfig(cfg)

		cfg.Namespace = "default"
		cfg.Context = context2

		require.NoError(t, err)
		assert.Equal(t, cfg, obj.Config())
	})

	t.Run("load with own context override namespace", func(t *testing.T) {
		t.Setenv("TestNewWithConfigSuccess."+t.Name(), "no-parallel")

		cfg := k8s.Config{Name: "name", Namespace: "own", Type: k8s.Secret, Context: context3}

		obj, err := k8s.NewWithConfig(cfg)

		cfg.Namespace = "own"
		cfg.Context = context3

		require.NoError(t, err)
		assert.Equal(t, cfg, obj.Config())
	})

	t.Run("load with own context override default namespace", func(t *testing.T) {
		t.Setenv("TestNewWithConfigSuccess."+t.Name(), "no-parallel")

		cfg := k8s.Config{Name: "name", Namespace: "self", Type: k8s.Secret, Context: context2}

		obj, err := k8s.NewWithConfig(cfg)

		cfg.Namespace = "self"
		cfg.Context = context2

		require.NoError(t, err)
		assert.Equal(t, cfg, obj.Config())
	})

	t.Run("load as current context", func(t *testing.T) {
		t.Setenv("TestNewWithConfigSuccess."+t.Name(), "no-parallel")

		cfg := k8s.Config{Name: "name", Namespace: "", Type: k8s.Secret, Context: context1}

		obj, err := k8s.NewWithConfig(cfg)

		cfg.Namespace = "myspace"
		cfg.Context = context1

		require.NoError(t, err)
		assert.Equal(t, cfg, obj.Config())
	})
}

func TestSourceExtract(t *testing.T) {
	t.Setenv("TestNewWithConfigSuccess", "no-parallel")
	t.Setenv("KUBECONFIG", testFile(t, successKonfig))

	const (
		testNamespace     = "test-namespace"
		testConfigMapName = "test-configmap"
		testSecretName    = "test-secret"
	)

	var (
		configMap = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: testConfigMapName, Namespace: testNamespace},
			Data:       map[string]string{"KEY1": "VALUE1", "KEY2": "VALUE2"},
		}

		secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: testSecretName, Namespace: testNamespace},
			Data:       map[string][]byte{"API_KEY": []byte("secret-value")},
		}
	)

	fakeClient := fake.NewClientset(configMap, secret).CoreV1()

	t.Run("extract from configmap", func(t *testing.T) {
		t.Setenv("TestSourceExtract."+t.Name(), "no-parallel")

		obj, err := k8s.NewWithConfig(
			k8s.Config{Name: testConfigMapName, Namespace: testNamespace, Type: k8s.ConfigMap, Context: ""},
		)

		require.NoError(t, err)

		got, err := obj.WithMocks(fakeClient).Extract(t.Context())

		require.NoError(t, err)
		assert.Equal(t, configMap.Data, got)
	})

	t.Run("configmap not found", func(t *testing.T) {
		t.Setenv("TestSourceExtract."+t.Name(), "no-parallel")

		obj, err := k8s.NewWithConfig(
			k8s.Config{Name: testConfigMapName, Namespace: "", Type: k8s.ConfigMap, Context: ""},
		)

		require.NoError(t, err)

		got, err := obj.WithMocks(fakeClient).Extract(t.Context())

		require.ErrorIs(t, err, k8s.ErrKubectlError)
		require.ErrorContains(t, err, `configmaps "test-configmap" not found`)
		assert.Nil(t, got)
	})

	t.Run("extract from secret", func(t *testing.T) {
		t.Setenv("TestSourceExtract."+t.Name(), "no-parallel")

		obj, err := k8s.NewWithConfig(
			k8s.Config{Name: testSecretName, Namespace: testNamespace, Type: k8s.Secret, Context: ""},
		)

		require.NoError(t, err)

		got, err := obj.WithMocks(fakeClient).Extract(t.Context())

		want := make(map[string]string)
		for key, val := range secret.Data {
			want[key] = string(val)
		}

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("secret not found", func(t *testing.T) {
		t.Setenv("TestSourceExtract."+t.Name(), "no-parallel")

		obj, err := k8s.NewWithConfig(k8s.Config{Name: testSecretName, Namespace: "", Type: k8s.Secret, Context: ""})

		require.NoError(t, err)

		got, err := obj.WithMocks(fakeClient).Extract(t.Context())

		require.ErrorIs(t, err, k8s.ErrKubectlError)
		require.ErrorContains(t, err, `secrets "test-secret" not found`)
		assert.Nil(t, got)
	})
}

func TestSourceAvailableContexts(t *testing.T) {
	t.Setenv("TestSourceAvailableContexts", "no-parallel")
	t.Setenv("KUBECONFIG", testFile(t, successKonfig))

	obj, err := k8s.NewWithConfig(k8s.Config{Name: "name", Namespace: "", Type: k8s.ConfigMap, Context: ""})

	require.NoError(t, err)

	got := obj.AvailableContexts()
	want := []string{context1, context2, context3}

	assert.Equal(t, want, got)
}
