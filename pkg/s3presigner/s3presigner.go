// Package s3presigner provides an S3 presigned URL generator for file uploads.
// It implements the livetemplate.Presigner interface for direct client-to-S3 uploads.
package s3presigner

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/livetemplate/livetemplate"
)

// S3Config configures AWS S3 presigned upload behavior.
type S3Config struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string        // Custom endpoint for MinIO/LocalStack
	Expiry          time.Duration // Default: 15 minutes
	KeyPrefix       string        // Optional prefix prepended to S3 key
}

// S3Presigner generates presigned PUT URLs for S3 uploads.
// It implements the livetemplate.Presigner interface.
type S3Presigner struct {
	client *s3.PresignClient
	config S3Config
}

// NewS3Presigner creates a new S3 presigner with the given configuration.
func NewS3Presigner(cfg S3Config) (*S3Presigner, error) {
	ctx := context.Background()

	if cfg.Expiry == 0 {
		cfg.Expiry = 15 * time.Minute
	}

	var awsCfg aws.Config
	var err error

	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			)),
		)
	} else {
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(cfg.Region),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	if cfg.Endpoint != "" {
		awsCfg.BaseEndpoint = aws.String(cfg.Endpoint)
	}

	s3Client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(s3Client)

	return &S3Presigner{
		client: presignClient,
		config: cfg,
	}, nil
}

// Presign generates a presigned PUT URL for uploading directly to S3.
func (p *S3Presigner) Presign(entry *livetemplate.UploadEntry) (livetemplate.UploadMeta, error) {
	key := p.generateKey(entry)

	req, err := p.client.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(p.config.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(entry.ClientType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = p.config.Expiry
	})

	if err != nil {
		return livetemplate.UploadMeta{}, fmt.Errorf("failed to presign S3 URL: %w", err)
	}

	return livetemplate.UploadMeta{
		Uploader: "s3",
		URL:      req.URL,
		Fields:   nil,
		Headers: map[string]string{
			"Content-Type": entry.ClientType,
		},
	}, nil
}

// generateKey creates S3 object key from upload entry.
// Format: {KeyPrefix}/{entryID}/{sanitized_filename}
func (p *S3Presigner) generateKey(entry *livetemplate.UploadEntry) string {
	filename := filepath.Base(entry.ClientName)

	if p.config.KeyPrefix != "" {
		return fmt.Sprintf("%s/%s/%s", p.config.KeyPrefix, entry.ID, filename)
	}
	return fmt.Sprintf("%s/%s", entry.ID, filename)
}
