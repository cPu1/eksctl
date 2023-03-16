package addon

import (
	_ "embed"
	"encoding/json"
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// AddonConfig holds the configuration for an addon.
type AddonConfig struct {
	Name                    string
	AWSManaged              bool
	SkipReadinessCheck      bool
	CoreAddon               bool
	KnownServiceAccountMeta *api.ClusterIAMMeta
	SetRecommendedPolicies  func(*api.Addon, clusterInfo) error
}

type clusterInfo interface {
	IPv6Enabled() bool
	GetRegion() string
}

//go:embed assets/cilium-operator-policy.json
var ciliumOperatorPolicyJSON []byte

var addonsConfig = []AddonConfig{
	{
		Name:       api.VPCCNIAddon,
		AWSManaged: true,
		CoreAddon:  true,
		KnownServiceAccountMeta: &api.ClusterIAMMeta{
			Name:      "aws-node",
			Namespace: "kube-system",
		},
		SetRecommendedPolicies: func(a *api.Addon, c clusterInfo) error {
			partition := api.Partition(c.GetRegion())
			if c.IPv6Enabled() {
				a.AttachPolicy = makeIPv6VPCCNIPolicyDocument(partition)
			} else {
				a.AttachPolicyARNs = []string{fmt.Sprintf("arn:%s:iam::aws:policy/%s", partition, api.IAMPolicyAmazonEKSCNIPolicy)}
			}
			return nil
		},
	},
	{
		Name:               api.CoreDNSAddon,
		CoreAddon:          true,
		AWSManaged:         true,
		SkipReadinessCheck: true,
	},
	{
		Name:       api.KubeProxyAddon,
		CoreAddon:  true,
		AWSManaged: true,
	},
	{
		Name: api.CiliumAddon,
		SetRecommendedPolicies: func(a *api.Addon, c clusterInfo) error {
			var policyDoc api.InlineDocument
			if err := json.Unmarshal(ciliumOperatorPolicyJSON, &policyDoc); err != nil {
				return fmt.Errorf("unexpected error decoding policy document: %w", err)
			}
			a.AttachPolicy = policyDoc
			return nil
		},
	},
	{
		Name:               api.AWSEBSCSIDriverAddon,
		AWSManaged:         true,
		SkipReadinessCheck: true,
		SetRecommendedPolicies: func(addon *api.Addon, _ clusterInfo) error {
			addon.WellKnownPolicies = api.WellKnownPolicies{
				EBSCSIController: true,
			}
			return nil
		},
	},
}
