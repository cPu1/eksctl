package vpc

import (
	"context"
	"fmt"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/weaveworks/eksctl/pkg/awsapi"
)

type ZoneType int

const (
	ZoneTypeAvailabilityZone ZoneType = iota
	ZoneTypeLocalZone
)

func DiscoverZoneTypes(ctx context.Context, ec2API awsapi.EC2, region string) (map[string]ZoneType, error) {
	output, err := ec2API.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("region-name"),
				Values: []string{region},
			}, {
				Name:   aws.String("state"),
				Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error describing zones: %w", err)
	}

	zoneTypeMapping := map[string]ZoneType{}
	for _, z := range output.AvailabilityZones {
		switch *z.ZoneType {
		case "availability-zone":
			zoneTypeMapping[*z.ZoneName] = ZoneTypeAvailabilityZone
		case "local-zone":
			zoneTypeMapping[*z.ZoneName] = ZoneTypeLocalZone
		default:
			return nil, fmt.Errorf("expected zone type to be local or AZ; got %g", *z.ZoneType)
		}
	}
	return zoneTypeMapping, nil
}
