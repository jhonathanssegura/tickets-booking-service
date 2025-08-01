package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client     *s3.Client
	BucketName string
}

func (s *S3Client) UploadTicketFile(ctx context.Context, key string, body io.Reader) error {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   body,
	})
	if err != nil {
		// Proporcionar mensajes de error más específicos
		var errorMsg string
		switch {
		case strings.Contains(err.Error(), "NoSuchBucket"):
			errorMsg = fmt.Sprintf("El bucket S3 '%s' no existe. Verifique que LocalStack esté ejecutándose y el bucket haya sido creado.", s.BucketName)
		case strings.Contains(err.Error(), "RequestCanceled"):
			errorMsg = "Error de conexión con S3. Verifique que LocalStack esté ejecutándose en http://localhost:4566."
		case strings.Contains(err.Error(), "AccessDenied"):
			errorMsg = fmt.Sprintf("Acceso denegado al bucket S3 '%s'. Verifique los permisos.", s.BucketName)
		default:
			errorMsg = fmt.Sprintf("Error subiendo archivo a S3: %v", err)
		}
		return fmt.Errorf(errorMsg)
	}
	return nil
}

func (s *S3Client) DownloadTicketFile(ctx context.Context, key string) (io.ReadCloser, error) {
	resp, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		var errorMsg string
		switch {
		case strings.Contains(err.Error(), "NoSuchBucket"):
			errorMsg = fmt.Sprintf("El bucket S3 '%s' no existe.", s.BucketName)
		case strings.Contains(err.Error(), "NoSuchKey"):
			errorMsg = fmt.Sprintf("El archivo '%s' no existe en el bucket S3 '%s'.", key, s.BucketName)
		case strings.Contains(err.Error(), "RequestCanceled"):
			errorMsg = "Error de conexión con S3. Verifique que LocalStack esté ejecutándose."
		default:
			errorMsg = fmt.Sprintf("Error descargando archivo de S3: %v", err)
		}
		return nil, fmt.Errorf(errorMsg)
	}
	return resp.Body, nil
}

// EnsureBucketExists verifica que el bucket existe y lo crea si es necesario
func (s *S3Client) EnsureBucketExists(ctx context.Context) error {
	_, err := s.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.BucketName),
	})
	if err != nil {
		// Intentar crear el bucket
		_, err = s.Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(s.BucketName),
		})
		if err != nil {
			return fmt.Errorf("error creando bucket S3 '%s': %v", s.BucketName, err)
		}
	}
	return nil
}
