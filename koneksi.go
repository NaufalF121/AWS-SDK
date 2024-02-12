package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type B2 interface {
	List() ([]string, error)
	Upload(string, string) error
	Download(string, string) error
	Delete(string) error
}
type B2Client struct {
	bucketName string
	s3Client   *s3.S3
}

func NewClient(endpoint string, keyID string, applicationKey string, bucketID string, region string, token string) (B2, error) {
	cfg := &aws.Config{
		Region:      &region,
		Endpoint:    &endpoint,
		Credentials: credentials.NewStaticCredentials(keyID, applicationKey, token),
	}
	Newsession, err := session.NewSession(cfg)
	if err != nil {
		panic(err)
	}

	s3Client := s3.New(Newsession)

	return &B2Client{
		bucketID,
		s3Client,
	}, nil
}

func (b *B2Client) List() ([]string, error) {
	objects, err := b.s3Client.ListObjects(&s3.ListObjectsInput{
		Bucket: &b.bucketName,
	})
	if err != nil {
		return nil, fmt.Errorf("[err][b2] failed to list objects: '%s", err)
	}

	result := make([]string, 0, len(objects.Contents))
	for _, obj := range objects.Contents {
		result = append(result, *obj.Key)
	}

	return result, nil
}

func (b *B2Client) Delete(fileName string) error {
	_, err := b.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &b.bucketName,
		Key:    &fileName,
	})
	return err
}

func (b *B2Client) Upload(fileName string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("[err][b2][upload] failed to open file: '%s'", err)
	}
	defer file.Close()

	uploader := s3manager.NewUploaderWithClient(b.s3Client)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: &b.bucketName,
		Key:    &fileName,
		Body:   file,
	})

	if err != nil {
		return fmt.Errorf("[err][b2][upload] failed to upload file: '%s'", err)
	}
	return nil
}

func (b *B2Client) Download(fileName string, writeTo string) error {
	file, err := os.Create(writeTo)
	if err != nil {
		return fmt.Errorf("[err][b2][download] err creating destination file: '%s'", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("[err][b2][download] failed to close file: '%s'", err)
		}
	}(file)

	downloader := s3manager.NewDownloaderWithClient(b.s3Client)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: &b.bucketName,
		Key:    &fileName,
	})

	if err != nil {
		return fmt.Errorf("[err][b2][download] failed to download file: '%s'", err)
	}
	return nil
}
