package providers

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (a *AwsProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	tm := transfermanager.New(a.GetS3Client())
	result, err := tm.UploadObject(context.TODO(), &transfermanager.UploadObjectInput{
		Bucket: aws.String(a.S3Bucket),
		Key:    aws.String(path),
		Body:   src,
	})

	if err != nil {
		return "", err
	}

	return a.publicURL(*result.Key), nil
}

func (a *AwsProvider) publicURL(key string) string {
	if a.s3Endpoint != "" {
		return fmt.Sprintf("%s/%s/%s", a.s3Endpoint, a.S3Bucket, key)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", a.S3Bucket, a.client.Region, key)
}

func (a *AwsProvider) DeleteFile(path string) error {
	_, err := a.GetS3Client().DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(a.S3Bucket),
		Key:    aws.String(strings.TrimPrefix(path, "/")),
	})

	return err
}
