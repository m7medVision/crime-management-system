package util

import (
	"context"
	"fmt"

	"github.com/m7medVision/crime-management-system/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

// InitMinio initializes the MinIO client
func InitMinio(cfg *config.Config) error {
	client, err := minio.New(cfg.Storage.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Storage.Minio.AccessKey, cfg.Storage.Minio.SecretKey, ""),
		Secure: cfg.Storage.Minio.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	// Check if the bucket exists or create it
	exists, err := client.BucketExists(context.Background(), cfg.Storage.Minio.Bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = client.MakeBucket(context.Background(), cfg.Storage.Minio.Bucket, minio.MakeBucketOptions{Region: cfg.Storage.Minio.Region})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	minioClient = client
	fmt.Println("MinIO client initialized successfully")
	return nil
}

// GetMinioClient returns the initialized MinIO client
func GetMinioClient() *minio.Client {
	return minioClient
}
