package k8s

import (
	"context"

	"github.com/therenotomorrow/ex"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	defaultNamespace = "default"

	ConfigMap ResourceType = "ConfigMap"
	Secret    ResourceType = "Secret"

	ErrMissingName  ex.C = "missing name"
	ErrMissingType  ex.C = "missing type"
	ErrInvalidType  ex.C = "invalid type"
	ErrKubectlError ex.C = "kubelib error"
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
		kubeRC clientcmdapi.Config
		client *kubernetes.Clientset
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

func (c *Config) merge(kubeRC *clientcmdapi.Config) {
	if c.Context == "" {
		c.Context = kubeRC.CurrentContext
	}

	if cc := kubeRC.Contexts[c.Context]; cc != nil {
		c.Namespace = cc.Namespace
	}

	if c.Namespace == "" {
		c.Namespace = defaultNamespace
	}
}

func NewWithConfig(config Config) (*Source, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	var overrides clientcmd.ConfigOverrides

	overrides.CurrentContext = config.Context

	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(), &overrides,
	)

	kubeRC, err := loader.RawConfig()
	if err != nil {
		return nil, ErrKubectlError.Because(err)
	}

	config.merge(&kubeRC)

	kubeConfig, err := loader.ClientConfig()
	if err != nil {
		return nil, ErrKubectlError.Because(err)
	}

	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, ErrKubectlError.Because(err)
	}

	return &Source{config: config, kubeRC: kubeRC, client: client}, nil
}

func (s *Source) Config() Config {
	return s.config
}

func (s *Source) AvailableContexts() []string {
	contexts := make([]string, 0, len(s.kubeRC.Contexts))
	for name := range s.kubeRC.Contexts {
		contexts = append(contexts, name)
	}

	return contexts
}

func (s *Source) Extract(ctx context.Context) (map[string]string, error) {
	var (
		options metav1.GetOptions
		mapping map[string]string
	)

	switch s.config.Type {
	case ConfigMap:
		configMap, err := s.client.CoreV1().ConfigMaps(s.config.Namespace).Get(ctx, s.config.Name, options)
		if err != nil {
			return nil, ErrKubectlError.Because(err)
		}

		mapping = configMap.Data

	case Secret:
		secret, err := s.client.CoreV1().Secrets(s.config.Namespace).Get(ctx, s.config.Name, options)
		if err != nil {
			return nil, ErrKubectlError.Because(err)
		}

		mapping = make(map[string]string)
		for key, val := range secret.Data {
			mapping[key] = string(val)
		}
	}

	return mapping, nil
}
