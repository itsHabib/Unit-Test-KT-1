package storage

//go:generate ...

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Service struct {
	s3         *s3.S3
	downloader *s3manager.Downloader
}

type Storage interface {
	Upload()
	Download()
}

func NewService(bucket string) (*Service, error) {
	if bucket == "" {
		return nil, fmt.Errorf("bucket is required")
	}

	// upload image to s3
	sess, err := session.NewSession(aws.NewConfig().WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to create aws session: %v", err)
	}

	return &Service{
		s3:         s3.New(sess),
		downloader: s3manager.NewDownloader(sess),
	}, nil
}

func (s *Service) Upload(body []byte, key string) error {
	if _, err := s.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("go-kt"),
		Key:    &key,
		Body:   bytes.NewReader(body),
	}); err != nil {
		return fmt.Errorf("unable to upload cat image: %w", err)
	}

	return nil
}

func (s *Service) Download(file, key string) error {
	// download image to file
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("unable to create file: %v", err)
	}
	if _, err := s.downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String("go-kt"),
		Key:    &key,
	}); err != nil {
		return fmt.Errorf("unable to download cat image: %w", err)
	}

	return nil
}
