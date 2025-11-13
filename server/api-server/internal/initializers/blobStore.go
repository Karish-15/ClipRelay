package initializers

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMinIO() *minio.Client {
	client, err := minio.New(
		os.Getenv("MINIO_ENDPOINT"),
		&minio.Options{
			Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
			Secure: true,
			Region: os.Getenv("MINIO_REGION"),
		},
	)
	if err != nil {
		log.Fatalf("failed to initialize Backblaze B2: %v", err)
	}

	return client
}

func InitializeMinIOBucket(client *minio.Client) {
	ctx := context.Background()
	bucket := os.Getenv("MINIO_BUCKET")

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil || !exists {
		log.Fatalf("bucket %s does not exist", bucket)
	}
}

func CreateAndInitMinIO() *minio.Client {
	client := ConnectMinIO()
	InitializeMinIOBucket(client)
	return client
}
