package builder

import (
	"strings"

	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

const (
	VPCResourceKey, IGWKey, GAKey                          = "VPC", "InternetGateway", "VPCGatewayAttachment"
	IPv6CIDRBlockKey                                       = "IPv6CidrBlock"
	EgressOnlyInternetGatewayKey                           = "EgressOnlyInternetGateway"
	ElasticIPKey                                           = "EIP"
	InternetCIDR, InternetIPv6CIDR                         = "0.0.0.0/0", "::/0"
	PubRouteTableKey, PrivateRouteTableKey                 = "PublicRouteTable", "PrivateRouteTable"
	PubRouteTableAssociation, PrivateRouteTableAssociation = "RouteTableAssociationPublic", "RouteTableAssociationPrivate"
	PubSubRouteKey, PubSubIPv6RouteKey                     = "PublicSubnetDefaultRoute", "PublicSubnetIPv6DefaultRoute"
	PrivateSubnetRouteKey, PrivateSubnetIpv6RouteKey       = "PrivateSubnetDefaultRoute", "PrivateSubnetDefaultIpv6Route"
	PublicSubnetKey, PrivateSubnetKey                      = "PublicSubnet", "PrivateSubnet"
	NATGatewayKey                                          = "NATGateway"
	PublicSubnetsOutputKey, PrivateSubnetsOutputKey        = "SubnetsPublic", "SubnetsPrivate"
	// AzA, AzB                             = "us-west-2a", "us-west-2b"
	// PrivateSubnet1, PrivateSubnet2       = "subnet-0ade11bad78dced9f", "subnet-0f98135715dfcf55a"
	// publicSubnet1, publicSubnet2         = "subnet-0ade11bad78dced9e", "subnet-0f98135715dfcf55f"
	// privateSubnetRef1, privateSubnetRef2 = "SubnetPrivateUSWEST2A", "SubnetPrivateUSWEST2B"
	// publicSubnetRef1, publicSubnetRef2   = "SubnetPublicUSWEST2A", "SubnetPublicUSWEST2B"
	// privRouteTableA, privRouteTableB     = "PrivateRouteTableUSWEST2A", "PrivateRouteTableUSWEST2B"
	// rtaPublicA, rtaPublicB               = "RouteTableAssociationPublicUSWEST2A", "RouteTableAssociationPublicUSWEST2B"
	// rtaPrivateA, rtaPrivateB             = "RouteTableAssociationPrivateUSWEST2A", "RouteTableAssociationPrivateUSWEST2B"
)

func formatAZ(az string) string {
	return strings.ToUpper(strings.ReplaceAll(az, "-", ""))
}

func getSubnetIPv6CIDRBlock() *gfnt.Value {
	// get 8 of /64 subnets from the auto-allocated IPv6 block,
	// and pick one block based on subnetIndexForIPv6 counter;
	// NOTE: this is done inside of CloudFormation using Fn::Cidr,
	// we don't slice it here, just construct the JSON expression
	// that does slicing at runtime.
	refIPv6CIDRv6 := gfnt.MakeFnSelect(
		gfnt.NewInteger(0), gfnt.MakeFnGetAttString("VPC", "Ipv6CidrBlocks"),
	)
	refSubnetSlices := gfnt.MakeFnCIDR(refIPv6CIDRv6, gfnt.NewInteger(8), gfnt.NewInteger(64))
	return refSubnetSlices
}

func getSubnetIPv4CIDRBlock() *gfnt.Value {
	refSubnetSlices := gfnt.MakeFnCIDR(gfnt.MakeFnGetAttString("VPC", "CidrBlocks"), gfnt.NewInteger(4), gfnt.NewInteger(14))
	return refSubnetSlices
}
