package maleos3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tigorlazuardi/maleo/bucket"
	"strings"
	"time"
)

type S3 struct {
	*s3.Client
	bucket     string
	region     string
	endpoint   string
	secure     bool
	urlBuilder URLBuilder
	pretext    fmt.Stringer
	resolver   s3.EndpointResolver
	isAws      bool
	uploader   *manager.Uploader

	// fallback should only be used on NewS3Bucket call.
	fallbackClient func() (*s3.Client, error)
}

func (s *S3) ensureBucket() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, _ = s.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(s.bucket),
	}, func(options *s3.Options) {
		options.Region = s.region
		if s.resolver != nil {
			options.EndpointResolver = s.resolver
		}
	})
}

// NewS3Bucket creates new S3 bucket. endpoint is basically the host of the bucket.
// e.g. mybucket.s3.amazonaws.com or localhost:9000.
// Do not include the protocol (http/https) in the endpoint.
// If the endpoint is empty, it will default to "s3.amazonaws.com".
//
// The default options assumes the machine that is running contains AWS credentials.
//
// If you wish to use a different or customized client or a client with custom logic, use WithClient option.
//
// If your endpoint is like this:
//
//	<bucket>.s3.<region>.amazonaws.com
//
// the bucket and region will be detected automatically.
//
// If your endpoint is like this:
//
//	s3.<region>.amazonaws.com
//
// the region will be detected automatically, but you must
// provide the bucket name via WithBucket option.
//
// Non AWS endpoint will be treated as a custom endpoint, and will be treated as immutable hostname by the SDK.
// Also, you must provide the region via WithRegion option and the bucket name via WithBucket option.
//
// You may also want to modify how the URL is built when using custom endpoints, use WithURLBuilder option to change those.
func NewS3Bucket(endpoint string, opts ...Option) (bucket.Bucket, error) {
	if endpoint == "" {
		endpoint = "s3.amazonaws.com"
	}
	isAwsEndpoint := isAws(endpoint)
	bkt := detectBucket(endpoint)
	s := &S3{
		endpoint:   endpoint,
		secure:     true,
		region:     detectRegion(endpoint),
		bucket:     bkt,
		urlBuilder: s3URLBuilder(isAwsEndpoint),
		pretext:    StringerFunc(func() string { return "" }),
		isAws:      isAwsEndpoint,
	}
	for _, opt := range opts {
		opt.apply(s)
	}
	if s.bucket == "" {
		return nil, fmt.Errorf("bucket name is unknown, please provide it via WithBucket option")
	}
	if isAwsEndpoint && s.region == "" {
		return nil, fmt.Errorf("region is unknown, please provide it via WithRegion option")
	}
	if isAwsEndpoint {
		s.endpoint = fmt.Sprintf("%s.s3.%s.amazonaws.com", s.bucket, s.region)
	}
	if s.Client == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 15)
		defer cancel()
		defaultConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(s.region))
		if err != nil {
			return nil, err
		}
		s.Client = s3.NewFromConfig(defaultConfig)
	}
	if !s.isAws {
		s.resolver = s3.EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
			proto := "https"
			if !s.secure {
				proto = "http"
			}
			return aws.Endpoint{
				URL:               fmt.Sprintf("%s://%s", proto, s.endpoint),
				SigningRegion:     s.region,
				HostnameImmutable: true,
			}, nil
		})
	}
	s.uploader = manager.NewUploader(s.Client)
	s.ensureBucket()
	return s, nil
}

func isAws(endpoint string) bool {
	if endpoint == "" {
		return true
	}
	return strings.Contains(endpoint, "amazonaws.com")
}

func detectRegion(endpoint string) string {
	split := strings.Split(endpoint, ".")
	if len(split) < 4 {
		return ""
	}
	// handle s3.<region>.amazonaws.com
	if split[0] == "s3" && split[2] == "amazonaws" {
		return split[1]
	}
	if len(split) < 5 {
		return ""
	}
	// handle <bucket>.s3.<region>.amazonaws.com
	if split[1] == "s3" && split[3] == "amazonaws" {
		return split[2]
	}
	return ""
}

func detectBucket(endpoint string) string {
	split := strings.Split(endpoint, ".")
	if len(split) < 5 {
		return ""
	}
	// handle <bucket>.s3.<region>.amazonaws.com
	if split[1] == "s3" && split[3] == "amazonaws" {
		return split[0]
	}
	return ""
}

type SimpleStaticParams struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Endpoint        string
	Region          string
	Secure          bool
}

// SimpleStaticClient creates a new S3 client from passed credentials.
func SimpleStaticClient(params SimpleStaticParams, opts ...func(options *config.LoadOptions) error) (*s3.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	opts = append(opts, []func(options *config.LoadOptions) error{
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(params.AccessKeyID, params.SecretAccessKey, params.SessionToken),
		),
	}...)
	if !isAws(params.Endpoint) {
		opts = append(opts, createResolverOption(params.Endpoint, params.Region, params.Secure))
	}
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}

// createResolverOption creates an option that will set the endpoint resolver to a custom endpoint. Used for non AWS endpoints.
func createResolverOption(endpoint string, region string, secure bool) config.LoadOptionsFunc {
	return config.WithEndpointResolverWithOptions(
		aws.EndpointResolverWithOptionsFunc(
			func(service, _region string, options ...interface{}) (aws.Endpoint, error) {
				proto := "https"
				if !secure {
					proto = "http"
				}
				return aws.Endpoint{
					URL:               fmt.Sprintf("%s://%s", proto, endpoint),
					SigningRegion:     region,
					HostnameImmutable: true,
				}, nil
			},
		),
	)
}
