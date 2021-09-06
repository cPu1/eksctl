package testutils

import (
	"bytes"
	"encoding/json"
	"io"

	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

func ClusterConfigReader(clusterConfig *api.ClusterConfig) io.Reader {
	data, err := json.Marshal(clusterConfig)
	Expect(err).ToNot(HaveOccurred())
	return bytes.NewReader(data)
}

func MakeRawClient(clusterName, region string) *kubernetes.RawClient {
	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   clusterName,
			Region: region,
		},
	}
	ctl, err := eks.New(&api.ProviderConfig{Region: region}, cfg)
	Expect(err).NotTo(HaveOccurred())

	err = ctl.RefreshClusterStatus(cfg)
	Expect(err).ShouldNot(HaveOccurred())
	rawClient, err := ctl.NewRawClient(cfg)
	Expect(err).ToNot(HaveOccurred())
	return rawClient
}
