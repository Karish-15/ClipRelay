package initializers

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMinIO() *minio.Client {
	internalEndpoint := os.Getenv("MINIO_ENDPOINT")      // e.g. cliprelay-minio:9000
	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT") // e.g. 127.0.0.1:9000
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	useSSL, err := strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))
	if err != nil {
		useSSL = false
	}

	// Custom dialer: whenever MinIO tries to connect to "publicEndpoint",
	// actually connect to "internalEndpoint" (inside Docker network).
	customDial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		if addr == publicEndpoint {
			addr = internalEndpoint
		}
		d := net.Dialer{Timeout: 5 * time.Second}
		return d.DialContext(ctx, network, addr)
	}

	transport := &http.Transport{
		DialContext: customDial,
	}

	client, err := minio.New(publicEndpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:    useSSL,
		Transport: transport,
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
