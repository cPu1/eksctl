package vpc_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

type vpcImportEntry struct {
	expectedErr string
}

var _ = FDescribe("VPC Import", func() {
	DescribeTable("Import VPC and subnets", func(e vpcImportEntry) {
		provider := mockprovider.NewMockProvider()
		stack := &cloudformation.Stack{
			Outputs: []*cloudformation.Output{
				{
					OutputKey:   aws.String(outputs.ClusterVPC),
					OutputValue: aws.String("vpc-123"),
				},
				{
					OutputKey:   aws.String(outputs.ClusterSecurityGroup),
					OutputValue: aws.String("sg-123"),
				},
				{
					OutputKey:   aws.String(outputs.ClusterSubnetsPrivate),
					OutputValue: aws.String("subnet-pri1,subnet-pri2,subnet-pri3"),
				},
				{
					OutputKey:   aws.String(outputs.ClusterSubnetsPublic),
					OutputValue: aws.String("subnet-pub1,subnet-pub2,subnet-pub3"),
				},
				{
					OutputKey:   aws.String(outputs.ClusterSubnetsPrivateLocal),
					OutputValue: aws.String("subnet-l1,subnet-l2,subnet-l3"),
				},
				{
					OutputKey:   aws.String(outputs.ClusterSubnetsPublicLocal),
					OutputValue: aws.String("subnet-publ1,subnet-publ2,subnet-publ3"),
				},
			},
		}
		clusterConfig := &v1alpha5.ClusterConfig{
			Metadata: &v1alpha5.ClusterMeta{
				Name: "test",
			},
		}
		provider.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *eks.DescribeClusterInput) bool {
			return *input.Name == clusterConfig.Metadata.Name
		})).Return(&eks.DescribeClusterOutput{
			Cluster: &eks.Cluster{
				ResourcesVpcConfig: &eks.VpcConfigResponse{
					EndpointPrivateAccess: aws.Bool(false),
					EndpointPublicAccess:  aws.Bool(true),
				},
			},
		}, nil)

		provider.MockEC2().On("DescribeVpcs", mock.Anything, mock.Anything).Return(&ec2.DescribeVpcsOutput{
			Vpcs: []ec2types.Vpc{
				{
					CidrBlock: aws.String("192.168.0.0/19"),
					VpcId:     aws.String("vpc-123"),
					CidrBlockAssociationSet: []ec2types.VpcCidrBlockAssociation{
						{
							CidrBlock: aws.String("192.168.0.0/19"),
						},
					},
				},
			},
		}, nil)

		provider.MockEC2().On("DescribeSubnets", mock.Anything, mock.MatchedBy(func(input *ec2.DescribeSubnetsInput) bool {
			fmt.Println("desc", input.SubnetIds, input.Filters)
			return true
		})).Return(func(_ context.Context, input *ec2.DescribeSubnetsInput, _ ...func(*ec2.Options)) *ec2.DescribeSubnetsOutput {
			var subnets []ec2types.Subnet
			for i, s := range input.SubnetIds {
				subnets = append(subnets, ec2types.Subnet{
					SubnetId:         aws.String(s),
					AvailabilityZone: aws.String("us-west-2-" + fmt.Sprintf("%v", rune('a'+i))),
					VpcId:            aws.String("vpc-123"),
					CidrBlock:        aws.String("192.168.0.16/19"),
				})
			}
			return &ec2.DescribeSubnetsOutput{
				Subnets: subnets,
			}
		}, func(_ context.Context, _ *ec2.DescribeSubnetsInput, _ ...func(*ec2.Options)) error {
			return nil
		})

		err := vpc.UseFromClusterStack(context.Background(), provider, stack, clusterConfig)
		fmt.Println("test", clusterConfig)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}
		Expect(err).NotTo(HaveOccurred())
		fmt.Println("test", clusterConfig)
	},
		Entry("No subnets", vpcImportEntry{
			expectedErr: "",
		}),
	)
})
