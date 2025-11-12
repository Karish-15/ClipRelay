package initializers

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMinIO() *minio.Client {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	useSSL, err := strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))
	if err != nil {
		useSSL = false
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("failed to initialize MinIO client: %v", err)
	}

	return client
}

func InitializeMinIOBucket(client *minio.Client) {
	ctx := context.Background()
	bucket := os.Getenv("MINIO_BUCKET")

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatalf("failed to check bucket existence: %v", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("failed to create MinIO bucket: %v", err)
		}
	}
}

func CreateAndInitMinIO() *minio.Client {
	client := ConnectMinIO()
	InitializeMinIOBucket(client)
	return client
}
