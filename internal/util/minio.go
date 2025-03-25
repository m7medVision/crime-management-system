package util

import (
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

// InitMinio initializes the MinIO client
func InitMinio() error {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	// Check if the bucket exists or create it
	bucketName := os.Getenv("MINIO_BUCKET")
	location := os.Getenv("MINIO_REGION")
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
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
