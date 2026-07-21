package providers

import (
	"context"
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

	return *result.Key, nil
}

func (a *AwsProvider) DeleteFile(path string) error {
	_, err := a.GetS3Client().DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(a.S3Bucket),
		Key:    aws.String(strings.TrimPrefix(path, "/")),
	})

	return err
}
