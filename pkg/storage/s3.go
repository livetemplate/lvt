package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3StoreConfig configures an S3-backed file store.
type S3StoreConfig struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string // Custom endpoint for MinIO/LocalStack
	CDNPrefix       string // Optional CDN URL prefix (e.g., "https://cdn.example.com")
}

// S3Store stores files in AWS S3.
type S3Store struct {
	client *s3.Client
	config S3StoreConfig
}

// NewS3Store creates a new S3Store.
func NewS3Store(cfg S3StoreConfig) (*S3Store, error) {
	ctx := context.Background()

	var awsCfg aws.Config
	var err error

	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		awsCfg, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(cfg.Region),
			awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			)),
		)
	} else {
		awsCfg, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(cfg.Region),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("storage: load AWS config: %w", err)
	}

	if cfg.Endpoint != "" {
		awsCfg.BaseEndpoint = aws.String(cfg.Endpoint)
	}

	return &S3Store{
		client: s3.NewFromConfig(awsCfg),
		config: cfg,
	}, nil
}

// Save uploads the contents of reader to the given key in S3.
func (s *S3Store) Save(ctx context.Context, key string, reader io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("storage: s3 put: %w", err)
	}
	return nil
}

// Open returns a ReadCloser for the object at the given key in S3.
func (s *S3Store) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("storage: s3 get: %w", err)
	}
	return out.Body, nil
}

// Delete removes the object at the given key from S3.
func (s *S3Store) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("storage: s3 delete: %w", err)
	}
	return nil
}

// URL returns the public URL for the given key.
// Uses CDNPrefix if configured, otherwise constructs the standard S3 URL.
func (s *S3Store) URL(key string) string {
	if s.config.CDNPrefix != "" {
		return s.config.CDNPrefix + "/" + key
	}
	if s.config.Endpoint != "" {
		return fmt.Sprintf("%s/%s/%s", s.config.Endpoint, s.config.Bucket, key)
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.config.Bucket, s.config.Region, key)
}
