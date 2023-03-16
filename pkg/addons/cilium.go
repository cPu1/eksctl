package addons

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cilium/cilium-cli/status"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"helm.sh/helm/v3/pkg/cli/values"

	"github.com/kris-nova/logger"

	"k8s.io/client-go/dynamic"
	restclient "k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"

	ciliumDefaults "github.com/cilium/cilium-cli/defaults"
	ciliumInstall "github.com/cilium/cilium-cli/install"
	"github.com/cilium/cilium-cli/k8s"
	ciliumClientSet "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// KubernetesProvider provides an interface for constructing a Kubernetes client.
type KubernetesProvider interface {
	GetConfig() clientcmdapi.Config
	RawConfig() *restclient.Config
	NewClientSet() (*kubernetes.Clientset, error)
}

// ClusterInfo holds the cluster info.
type ClusterInfo interface {
	// Meta returns the cluster metadata.
	Meta() *api.ClusterMeta
	// GetRegion returns the region.
	GetRegion() string
}

// Cilium allows installing Cilium as an addon.
type Cilium struct {
	// KubernetesProvider allows interacting with Kubernetes.
	KubernetesProvider KubernetesProvider
	// ClusterInfo holds the cluster info.
	ClusterInfo ClusterInfo
	// ServiceAccountRoleARN configures the Cilium operator to use IAM Roles for Service Accounts.
	ServiceAccountRoleARN string
}

// Install installs the Cilium addon.
func (c *Cilium) Install(ctx context.Context, configValues map[string]interface{}) error {
	client, err := makeCiliumClient(c.KubernetesProvider)
	if err != nil {
		return err
	}

	if err := c.setDefaults(configValues); err != nil {
		return err
	}

	// TODO: Modify cilium-cli to accept Helm values as a map, instead of having to pass a values file.
	valuesFile, err := makeHelmValuesFile(configValues)
	if err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(valuesFile); err != nil {
			logger.Warning("error removing Helm values file: %v", err)
		}
	}()
	clusterMeta := c.ClusterInfo.Meta()
	params := ciliumInstall.Parameters{
		Namespace:      metav1.NamespaceSystem,
		DatapathMode:   "aws-eni",
		Writer:         os.Stdout,
		ListVersions:   false,
		NodeEncryption: false,
		ClusterID:      1,
		ClusterName:    clusterMeta.Name,
		// TODO: allow passing version.
		Version:              ciliumDefaults.Version,
		Encryption:           "disabled",
		HelmValuesSecretName: ciliumDefaults.HelmValuesSecretName,
		K8sVersion:           clusterMeta.Version,
		Rollback:             true,
		CiliumReadyTimeout:   5 * time.Minute,
		KubeProxyReplacement: "disabled",
		Wait:                 false,
		RestartUnmanagedPods: false,
		HelmOpts: values.Options{
			ValueFiles: []string{valuesFile},
			Values: []string{
				fmt.Sprintf("serviceAccounts.operator.annotations.%s=%s", strings.ReplaceAll(api.AnnotationEKSRoleARN, ".", `\.`), c.ServiceAccountRoleARN),
			},
		},
	}

	// TODO: remove this workaround.
	client.RawConfig.Contexts[""] = client.RawConfig.Contexts[client.RawConfig.CurrentContext]
	installer, err := ciliumInstall.NewK8sInstaller(client, params)

	if err != nil {
		return fmt.Errorf("error creating Cilium installer: %w", err)
	}
	// Mention that a lot of checks would be redundant in the CLI since it's ensured that are no nodes.
	return installer.Install(ctx)
}

func makeCiliumClient(provider KubernetesProvider) (*k8s.Client, error) {
	restConfig := provider.RawConfig()
	ciliumClientSet, err := ciliumClientSet.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating Cilium client: %w", err)
	}
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating dynamic client: %w", err)
	}
	clientset, err := provider.NewClientSet()
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes client: %w", err)
	}
	return &k8s.Client{
		Clientset:        clientset,
		DynamicClientset: dynamicClient,
		CiliumClientset:  ciliumClientSet,
		Config:           restConfig,
		RawConfig:        provider.GetConfig(),
	}, nil
}

// CiliumStatusCollector is an interface for querying the status of Cilium.
type CiliumStatusCollector interface {
	// Status reports the status of Cilium.
	Status(ctx context.Context) (*status.Status, error)
}

// NewCiliumStatusCollector creates and returns a new CiliumStatusCollector.
func NewCiliumStatusCollector(provider KubernetesProvider) (CiliumStatusCollector, error) {
	client, err := makeCiliumClient(provider)
	if err != nil {
		return nil, err
	}
	statusCollector, err := status.NewK8sStatusCollector(client, status.K8sStatusParameters{})
	if err != nil {
		return nil, fmt.Errorf("unexpected error creating Cilium status collector: %w", err)
	}
	return statusCollector, nil
}

func (c *Cilium) setDefaults(configValues map[string]interface{}) error {
	var operator map[string]interface{}
	const operatorKey = "operator"

	if op, ok := configValues[operatorKey]; !ok {
		operator = map[string]interface{}{}
		configValues[operatorKey] = operator
	} else if _, ok := op.(map[string]interface{}); !ok {
		return fmt.Errorf("expected %q to be a %T; got %T", operatorKey, operator, op)
	}

	operator["extraEnv"] = []corev1.EnvVar{
		{
			Name:  "AWS_REGION",
			Value: c.ClusterInfo.Meta().Region,
		},
	}
	return nil
}

func makeHelmValuesFile(config map[string]interface{}) (string, error) {
	helmValues, err := yaml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("error marshaling addon config: %w", err)
	}
	valuesFile, err := os.CreateTemp(os.TempDir(), "helm-values-*.yaml")
	if err != nil {
		return "", fmt.Errorf("error creating Helm values file: %w", err)
	}
	if _, err := valuesFile.Write(helmValues); err != nil {
		return "", fmt.Errorf("error writing to Helm values file: %w", err)
	}
	if err := valuesFile.Close(); err != nil {
		return "", fmt.Errorf("error closing Helm values file: %w", err)
	}
	return valuesFile.Name(), nil
}
