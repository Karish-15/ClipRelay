## Local MinIO Support (Custom Dialer for Internal ↔ Public Endpoint Mapping)

During local development, the backend (inside Docker) must connect to MinIO using the internal hostname:

```
cliprelay-minio:9000
```

but the browser can only access MinIO using the exposed host port:

```
127.0.0.1:9000
```

To solve this, ClipRelay uses a **custom dialer** so the MinIO SDK connects to the internal Docker hostname, while the client still receives signed URLs using the public host.  
This makes both uploads and downloads work seamlessly in local mode. Just Replace the function inside `api-server/internal/initializers/blobStore.go` and `bg-workers/mq-consumer/initializers/blob.go`

### Custom MinIO dialer (used inside Docker)

```go
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
```

This allows:

- Backend to **connect internally** to `cliprelay-minio:9000`
- Browser to load blobs via **public signed URLs** using `127.0.0.1:9000`
- Signatures remain valid (host is not part of S3 HMAC)
