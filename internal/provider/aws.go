package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ResourceFetcher defines the interface for fetching live cloud resource attributes.
type ResourceFetcher interface {
	FetchAttributes(ctx context.Context, resourceType, resourceID string) (map[string]interface{}, error)
}

// AWSProvider fetches live resource state from AWS.
type AWSProvider struct {
	s3Client *s3.Client
	region   string
}

// NewAWSProvider creates an AWSProvider using the default AWS credential chain.
func NewAWSProvider(ctx context.Context, region string) (*AWSProvider, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("loading aws config: %w", err)
	}
	return &AWSProvider{
		s3Client: s3.NewFromConfig(cfg),
		region:   region,
	}, nil
}

// FetchAttributes retrieves live attributes for the given resource type and ID.
func (p *AWSProvider) FetchAttributes(ctx context.Context, resourceType, resourceID string) (map[string]interface{}, error) {
	switch resourceType {
	case "aws_s3_bucket":
		return p.fetchS3Bucket(ctx, resourceID)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

func (p *AWSProvider) fetchS3Bucket(ctx context.Context, bucketName string) (map[string]interface{}, error) {
	out, err := p.s3Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return nil, fmt.Errorf("fetching s3 bucket %q: %w", bucketName, err)
	}
	return map[string]interface{}{
		"bucket": bucketName,
		"region": string(out.LocationConstraint),
	}, nil
}
