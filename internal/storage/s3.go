package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client envuelve el cliente de AWS S3 y el bucket a utilizar.
type S3Client struct {
	Client     *s3.Client
	BucketName string
}

// UploadTicketFile sube un archivo (por ejemplo, un PDF de ticket) a S3.
func (s *S3Client) UploadTicketFile(ctx context.Context, key string, body io.Reader) error {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   body,
	})
	if err != nil {
		return fmt.Errorf("error uploading file to S3: %w", err)
	}
	return nil
}

// DownloadTicketFile descarga un archivo desde S3 (por ejemplo, un PDF de ticket).
func (s *S3Client) DownloadTicketFile(ctx context.Context, key string) (io.ReadCloser, error) {
	resp, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("error downloading file from S3: %w", err)
	}
	return resp.Body, nil
}
