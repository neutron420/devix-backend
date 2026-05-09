package media

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"devix-backend/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type StorageProvider interface {
	Upload(ctx context.Context, file io.Reader, directory, filename string) (string, error)
	Delete(ctx context.Context, path string) error
	GetURL(path string) string
}

type LocalStorage struct {
	basePath string
	baseURL  string
}

func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

func (s *LocalStorage) Upload(_ context.Context, file io.Reader, directory, filename string) (string, error) {

	ext := filepath.Ext(filename)
	uniqueName := uuid.New().String() + ext

	dirPath := filepath.Join(s.basePath, directory)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(dirPath, uniqueName)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {

		os.Remove(filePath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	relativePath := filepath.Join(directory, uniqueName)
	return relativePath, nil
}

func (s *LocalStorage) Delete(_ context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *LocalStorage) GetURL(path string) string {
	return s.baseURL + "/" + filepath.ToSlash(path)
}

type R2Storage struct {
	client     *s3.Client
	bucketName string
	publicURL  string
	cdnURL     string
}

func NewR2Storage(ctx context.Context, cfg config.R2Config) (*R2Storage, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: cfg.Endpoint,
		}, nil
	})

	awsCfg, err := s3config.LoadDefaultConfig(ctx,
		s3config.WithEndpointResolverWithOptions(customResolver),
		s3config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		s3config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load R2 config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	return &R2Storage{
		client:     client,
		bucketName: cfg.BucketName,
		publicURL:  cfg.PublicURL,
		cdnURL:     cfg.CDNURL,
	}, nil
}

func (s *R2Storage) Upload(ctx context.Context, file io.Reader, directory, filename string) (string, error) {
	ext := filepath.Ext(filename)
	uniqueName := uuid.New().String() + ext
	key := filepath.ToSlash(filepath.Join(directory, uniqueName))

	_, err := s.client.PutObject(ctx, s.inputPutObject(key, file))
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %w", err)
	}

	return key, nil
}

func (s *R2Storage) inputPutObject(key string, file io.Reader) *s3.PutObjectInput {
	return &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
		Body:   file,
	}
}

func (s *R2Storage) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from R2: %w", err)
	}
	return nil
}

func (s *R2Storage) GetURL(path string) string {
	baseURL := s.publicURL
	if s.cdnURL != "" {
		baseURL = s.cdnURL
	}
	return baseURL + "/" + path
}
