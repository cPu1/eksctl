package builder_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

var _ = FDescribe("AssignSubnets", func() {
	type assignSubnetsEntry struct {
		np               api.NodePool
		mockEC2          func(provider *mockprovider.MockProvider)
		setSubnetMapping func(config *api.ClusterConfig)

		expectedErr       string
		expectedSubnetIDs []string
	}

	toSubnetIDs := func(subnetRefs *gfnt.Value) []string {
		subnetsSlice, ok := subnetRefs.Raw().(gfnt.Slice)
		Expect(ok).To(BeTrue(), fmt.Sprintf("expected subnet refs to be of type %T; got %T", gfnt.Slice{}, subnetRefs.Raw()))
		var subnetIDs []string
		for _, subnetID := range subnetsSlice {
			subnetIDs = append(subnetIDs, subnetID.String())
		}
		return subnetIDs
	}

	DescribeTable("assigns subnets to a nodegroup", func(e assignSubnetsEntry) {
		clusterConfig := api.NewClusterConfig()
		if e.setSubnetMapping != nil {
			e.setSubnetMapping(clusterConfig)
		}
		mockProvider := mockprovider.NewMockProvider()
		if e.mockEC2 != nil {
			e.mockEC2(mockProvider)
		}
		subnetRefs, err := builder.AssignSubnets(context.Background(), e.np, nil, clusterConfig, mockProvider.EC2())
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}
		Expect(err).NotTo(HaveOccurred())
		fmt.Println("subnets", subnetRefs)
		subnetIDs := toSubnetIDs(subnetRefs)
		Expect(err).NotTo(HaveOccurred())
		Expect(subnetIDs).To(ConsistOf(e.expectedSubnetIDs))

	},

		Entry("self-managed nodegroup with availability zones", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					AvailabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
				},
			},
			setSubnetMapping: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.VPC.Subnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-1a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-1a",
						},
						"us-west-1b": api.AZSubnetSpec{
							ID: "subnet-1b",
							AZ: "us-west-1b",
						},
						"us-west-1c": api.AZSubnetSpec{
							ID: "subnet-1c",
							AZ: "us-west-1c",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},
			expectedSubnetIDs: []string{"subnet-1a", "subnet-1b", "subnet-1c"},
		}),

		Entry("managed nodegroup with availability zones", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					AvailabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
				},
			},
			setSubnetMapping: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.VPC.Subnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-1a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-1a",
						},
						"us-west-1b": api.AZSubnetSpec{
							ID: "subnet-1b",
							AZ: "us-west-1b",
						},
						"us-west-1c": api.AZSubnetSpec{
							ID: "subnet-1c",
							AZ: "us-west-1c",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},
			expectedSubnetIDs: []string{"subnet-1a", "subnet-1b", "subnet-1c"},
		}),

		Entry("self-managed nodegroup with local zones", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{},
				LocalZones:    []string{"us-west-2-lax-1a", "us-west-2-lax-1b"},
			},
			setSubnetMapping: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.VPC.LocalZoneSubnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2-lax-1a": api.AZSubnetSpec{
							ID: "subnet-lax-1a",
							AZ: "us-west-2-lax-1a",
						},
						"us-west-2-lax-1b": api.AZSubnetSpec{
							ID: "subnet-lax-1b",
							AZ: "us-west-2-lax-1b",
						},
						"us-west-2-lax-1d": api.AZSubnetSpec{
							ID: "subnet-lax-1d",
							AZ: "us-west-2-lax-1d",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},

			expectedSubnetIDs: []string{"subnet-lax-1a", "subnet-lax-1b"},
			/*mockEC2: func(provider *mockprovider.MockProvider) {
				provider.MockEC2().On("DescribeLaunchTemplateVersions", mock.Anything, mock.MatchedBy(matcher)).
					Return(&ec2.DescribeLaunchTemplateVersionsOutput{
						LaunchTemplateVersions: []ec2types.LaunchTemplateVersion{
							{
								LaunchTemplateData: lt,
							},
						},
					}, nil)

			},*/
		}),

		Entry("self-managed nodegroup with privateNetworking and local zones", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					PrivateNetworking: true,
				},
				LocalZones: []string{"us-west-2-lax-1a", "us-west-2-lax-1b"},
			},
			setSubnetMapping: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.VPC.LocalZoneSubnets = &api.ClusterSubnets{
					Public: api.NewAZSubnetMapping(),
					Private: api.AZSubnetMapping{
						"us-west-2-lax-1a": api.AZSubnetSpec{
							ID: "subnet-lax-1a",
							AZ: "us-west-2-lax-1a",
						},
						"us-west-2-lax-1b": api.AZSubnetSpec{
							ID: "subnet-lax-1b",
							AZ: "us-west-2-lax-1b",
						},
						"us-west-2-lax-1d": api.AZSubnetSpec{
							ID: "subnet-lax-1d",
							AZ: "us-west-2-lax-1d",
						},
					},
				}
			},

			expectedSubnetIDs: []string{"subnet-lax-1a", "subnet-lax-1b"},
		}),

		FEntry("self-managed nodegroup with local zones and subnet IDs", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets: []string{"subnet-z1", "subnet-z2"},
				},
				LocalZones: []string{"us-west-2-lax-1a", "us-west-2-lax-1b"},
			},
			setSubnetMapping: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.VPC.LocalZoneSubnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2-lax-1a": api.AZSubnetSpec{
							ID: "subnet-lax-1a",
							AZ: "us-west-2-lax-1a",
						},
						"us-west-2-lax-1b": api.AZSubnetSpec{
							ID: "subnet-lax-1b",
							AZ: "us-west-2-lax-1b",
						},
						"us-west-2-lax-1d": api.AZSubnetSpec{
							ID: "subnet-lax-1d",
							AZ: "us-west-2-lax-1d",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},
			expectedSubnetIDs: []string{"subnet-z1", "subnet-z2", "subnet-lax-1a", "subnet-lax-1b"},

			mockEC2: func(provider *mockprovider.MockProvider) {
				provider.MockEC2().
					On("DescribeSubnets", mock.Anything, mock.Anything).Return(func(_ context.Context, input *ec2.DescribeSubnetsInput, _ ...func(options *ec2.Options)) *ec2.DescribeSubnetsOutput {
					return &ec2.DescribeSubnetsOutput{
						Subnets: []ec2types.Subnet{
							{
								SubnetId:         aws.String(input.SubnetIds[0]),
								AvailabilityZone: aws.String("us-west-2-lax-1e"),
								VpcId:            aws.String("vpc-1"),
							},
						},
					}
				}, nil).
					On("DescribeAvailabilityZones", mock.Anything, mock.Anything).Return(&ec2.DescribeAvailabilityZonesOutput{
					AvailabilityZones: []ec2types.AvailabilityZone{
						{
							ZoneType: aws.String("local-zone"),
							ZoneName: aws.String("us-west-2-lax-1e"),
						},
						{
							ZoneType: aws.String("availability-zone"),
							ZoneName: aws.String("us-west-2d"),
						},
						{
							ZoneType: aws.String("local-zone"),
							ZoneName: aws.String("us-west-2-lax-1f"),
						},
					},
				}, nil)

			},
		}),

		Entry("managed nodegroup with privateNetworking, availability zones and subnet IDs", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets:           []string{"subnet-z1", "subnet-z2"},
					AvailabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
				},
			},
			setSubnetMapping: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.VPC.Subnets = &api.ClusterSubnets{
					Private: api.AZSubnetMapping{
						"us-west-1a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-1a",
						},
						"us-west-1b": api.AZSubnetSpec{
							ID: "subnet-1b",
							AZ: "us-west-1b",
						},
						"us-west-1c": api.AZSubnetSpec{
							ID: "subnet-1c",
							AZ: "us-west-1c",
						},
					},
					Public: api.NewAZSubnetMapping(),
				}
			},
			expectedSubnetIDs: []string{"subnet-1a", "subnet-1b", "subnet-1c", "subnet-z1", "subnet-z2"},

			mockEC2: func(provider *mockprovider.MockProvider) {
				provider.MockEC2().
					On("DescribeSubnets", mock.Anything, mock.MatchedBy(func(_ context.Context, input *ec2.DescribeSubnetsInput) bool {
						return len(input.SubnetIds) == 2
					})).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []ec2types.Subnet{
						{
							SubnetId:         aws.String("subnet-z1"),
							AvailabilityZone: aws.String("us-west-2g"),
							VpcId:            aws.String("vpc-1"),
						},
						{
							SubnetId:         aws.String("subnet-z2"),
							AvailabilityZone: aws.String("us-west-2h"),
							VpcId:            aws.String("vpc-1"),
						},
					},
				}, nil).
					On("DescribeAvailabilityZones", mock.Anything, mock.Anything).Return(&ec2.DescribeAvailabilityZonesOutput{
					AvailabilityZones: []ec2types.AvailabilityZone{
						{
							ZoneType: aws.String("local-zone"),
							ZoneName: aws.String("us-west-2-lax-1e"),
						},
						{
							ZoneType: aws.String("availability-zone"),
							ZoneName: aws.String("us-west-2g"),
						},
						{
							ZoneType: aws.String("local-zone"),
							ZoneName: aws.String("us-west-2h"),
						},
					},
				}, nil)

			},
		}),

		Entry("managed nodegroup with availability zones and subnet IDs in local zones", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets:           []string{"subnet-z1", "subnet-z2"},
					AvailabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
				},
			},
			setSubnetMapping: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.VPC.Subnets = &api.ClusterSubnets{
					Private: api.AZSubnetMapping{
						"us-west-1a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-1a",
						},
						"us-west-1b": api.AZSubnetSpec{
							ID: "subnet-1b",
							AZ: "us-west-1b",
						},
						"us-west-1c": api.AZSubnetSpec{
							ID: "subnet-1c",
							AZ: "us-west-1c",
						},
					},
					Public: api.NewAZSubnetMapping(),
				}
			},

			expectedErr: "managed nodegroups cannot be launched in local zones",

			mockEC2: func(provider *mockprovider.MockProvider) {
				provider.MockEC2().
					On("DescribeSubnets", mock.Anything, mock.MatchedBy(func(input *ec2.DescribeSubnetsInput) bool {
						return len(input.SubnetIds) == 2
					})).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []ec2types.Subnet{
						{
							SubnetId:         aws.String("subnet-z1"),
							AvailabilityZone: aws.String("us-west-2-lax-1e"),
							VpcId:            aws.String("vpc-1"),
						},
						{
							SubnetId:         aws.String("subnet-z2"),
							AvailabilityZone: aws.String("us-west-2-lax-1f"),
							VpcId:            aws.String("vpc-1"),
						},
					},
				}, nil).
					On("DescribeAvailabilityZones", mock.Anything, mock.Anything).Return(ec2.DescribeAvailabilityZonesOutput{
					AvailabilityZones: []ec2types.AvailabilityZone{
						{
							ZoneType: aws.String("local-zone"),
							ZoneName: aws.String("us-west-2-lax-1e"),
						},
						{
							ZoneType: aws.String("availability-zone"),
							ZoneName: aws.String("us-west-2d"),
						},
						{
							ZoneType: aws.String("local-zone"),
							ZoneName: aws.String("us-west-2-lax-1f"),
						},
					},
				}, nil)

			},
		}),
	)
})
