package storage

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3UploaderInterface defines the interface for the S3 uploader
type S3UploaderInterface interface {
	UploadFile(folder, fileName, fileType, userID string, fileData []byte) (string, error)
	GeneratePresignedURL(key string, expiry time.Duration) (string, error)
}

// S3Uploader holds the S3 client and the bucket name
type S3Uploader struct {
	client *s3.Client
	bucket string
}

// NewS3Uploader creates an S3Uploader instance
func NewS3Uploader(accessKeyID, secretAccessKey, region, bucket string) (S3UploaderInterface, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Uploader{
		client: client,
		bucket: bucket,
	}, nil
}

func (u *S3Uploader) UploadFile(folder, fileName, fileType, userID string, fileData []byte) (string, error) {
	if fileName == "" {
		return "", fmt.Errorf("file name must not be empty")
	}

	// timestamp + fileName + userID
	now := time.Now().Unix()
	ext := filepath.Ext(fileName)
	base := strings.TrimSuffix(fileName, ext)
	safeBase := sanitizeFileName(base)
	finalFileName := fmt.Sprintf("%d_%s_%s%s", now, safeBase, userID, ext)

	// Build the S3 key: folder + finalFileName
	s3Key := finalFileName
	if folder != "" {
		s3Key = fmt.Sprintf("%s/%s", folder, finalFileName)
	}

	putInput := &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(fileData),
		ContentType: aws.String(fileType),
		ACL:         types.ObjectCannedACLPrivate,
	}

	// Execute the PUT request to S3
	_, err := u.client.PutObject(context.TODO(), putInput)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	return s3Key, nil
}

func (u *S3Uploader) GeneratePresignedURL(key string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(u.client, s3.WithPresignExpires(expiry))

	// Build the input
	getObjInput := &s3.GetObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
	}

	// Create the presigned GET URL
	presignResult, err := presignClient.PresignGetObject(context.TODO(), getObjInput)
	if err != nil {
		return "", fmt.Errorf("failed to presign GetObject: %v", err)
	}

	return presignResult.URL, nil
}

func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	return fmt.Sprintf("%s%s", base, ext)
}
