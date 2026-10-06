package minioclient

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client  *minio.Client
	bucket  string
	baseURL string
}

func New() (*Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucket)

	err = minioClient.SetBucketPolicy(ctx, bucket, policy)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:  minioClient,
		bucket:  bucket,
		baseURL: fmt.Sprintf("http://%s/%s", endpoint, bucket),
	}, nil
}

func (c *Client) UploadFile(header *multipart.FileHeader, prefix string) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	objectName := buildObjectName(header.Filename, prefix)

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = c.client.PutObject(
		context.Background(),
		c.bucket,
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (c *Client) FileURL(objectName string) string {
	if objectName == "" {
		return ""
	}

	if strings.HasPrefix(objectName, "http") || strings.HasPrefix(objectName, "/") {
		return objectName
	}

	return fmt.Sprintf("%s/%s", c.baseURL, objectName)
}

func buildObjectName(originalName string, prefix string) string {
	extension := strings.ToLower(filepath.Ext(originalName))

	clean := ""
	for _, symbol := range extension {
		if (symbol >= 'a' && symbol <= 'z') || (symbol >= '0' && symbol <= '9') || symbol == '.' {
			clean = clean + string(symbol)
		}
	}

	if clean == "" || clean == "." {
		clean = ".bin"
	}

	return fmt.Sprintf("%s_%d%s", prefix, time.Now().UnixNano(), clean)
}
