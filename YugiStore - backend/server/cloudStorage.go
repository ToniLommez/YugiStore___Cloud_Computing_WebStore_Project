package server

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
)

func UploadImage(ctx context.Context, file multipart.File, filename string) (string, error) {
	const bucketName = "yugistore-images"

	wc := StorageClient.Bucket(bucketName).Object(filename).NewWriter(ctx)
	wc.ContentType = "image/jpeg"
	// wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
	return publicURL, nil
}
