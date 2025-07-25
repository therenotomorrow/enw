package k8s

import (
	"context"
	"slices"

	"github.com/therenotomorrow/ex"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	defaultNamespace = "default"

	ConfigMap ResourceType = "configmap"
	Secret    ResourceType = "secret"

	ErrMissingName  ex.Const = "missing name"
	ErrMissingType  ex.Const = "missing type"
	ErrInvalidType  ex.Const = "invalid type"
	ErrKubectlError ex.Const = "kubelib error"
)

type (
	ResourceType string

	Config struct {
		Name      string
		Namespace string
		Type      ResourceType
		Context   string
	}

	Source struct {
		konfig api.Config
		client corev1.CoreV1Interface
		config Config
	}
)

func (c *Config) Validate() error {
	if c.Name == "" {
		return ErrMissingName
	}

	switch c.Type {
	case "":
		return ErrMissingType
	case ConfigMap, Secret:
	default:
		return ErrInvalidType
	}

	return nil
}

func (c *Config) merge(konfig *api.Config) {
	if c.Context == "" {
		c.Context = konfig.CurrentContext
	}

	if c.Namespace != "" {
		return
	}

	if cc := konfig.Contexts[c.Context]; cc != nil {
		c.Namespace = cc.Namespace
	}

	if c.Namespace == "" {
		c.Namespace = defaultNamespace
	}
}

func (c *Config) overrides() *clientcmd.ConfigOverrides {
	var overrides clientcmd.ConfigOverrides

	overrides.CurrentContext = c.Context

	return &overrides
}

func NewWithConfig(config Config) (*Source, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		config.overrides(),
	)

	konfig, err := loader.RawConfig()
	if err != nil {
		return nil, ErrKubectlError.Because(err)
	}

	config.merge(&konfig)

	clientKonfig, err := loader.ClientConfig()
	if err != nil {
		return nil, ErrKubectlError.Because(err)
	}

	client, err := kubernetes.NewForConfig(clientKonfig)
	if err != nil {
		return nil, ErrKubectlError.Because(err)
	}

	return &Source{config: config, konfig: konfig, client: client.CoreV1()}, nil
}

func (s *Source) Config() Config {
	return s.config
}

func (s *Source) WithMocks(mocks ...any) *Source {
	clone := &Source{config: s.config, konfig: s.konfig, client: s.client}

	for _, mock := range mocks {
		impl, ok := mock.(corev1.CoreV1Interface)
		if ok {
			clone.client = impl

			break
		}
	}

	return clone
}

func (s *Source) Extract(ctx context.Context) (map[string]string, error) {
	var (
		envs    map[string]string
		options metav1.GetOptions
	)

	switch s.config.Type {
	case ConfigMap:
		configMap, err := s.client.ConfigMaps(s.config.Namespace).Get(ctx, s.config.Name, options)
		if err != nil {
			return nil, ErrKubectlError.Because(err)
		}

		envs = configMap.Data

	case Secret:
		secret, err := s.client.Secrets(s.config.Namespace).Get(ctx, s.config.Name, options)
		if err != nil {
			return nil, ErrKubectlError.Because(err)
		}

		envs = make(map[string]string)
		for key, val := range secret.Data {
			envs[key] = string(val)
		}
	}

	return envs, nil
}

func (s *Source) AvailableContexts() []string {
	contexts := make([]string, 0, len(s.konfig.Contexts))
	for name := range s.konfig.Contexts {
		contexts = append(contexts, name)
	}

	slices.Sort(contexts)

	return contexts
}
