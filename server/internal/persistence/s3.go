package persistence

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Provider struct {
	client *s3.Client
	bucket string
}

type S3Config struct {
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
}

func NewS3Provider(cfg S3Config) (*S3Provider, error) {
	fmt.Println("creating new s3 provider")

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		awsconfig.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	return &S3Provider{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

func (s *S3Provider) Save(name string, data []byte) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &name,
		Body:   bytes.NewReader(data),
	})
	return err
}

func (s *S3Provider) Load(name string) ([]byte, error) {
	output, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &name,
	})

	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	var buf bytes.Buffer

	_, err = buf.ReadFrom(output.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
func (s *S3Provider) Delete(name string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &name,
	})
	return err
}
func (s *S3Provider) List() ([]string, error) {
	listObjectsOutput, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &s.bucket,
	})
	if err != nil {
		return nil, err
	}

	var objects []string

	for _, obj := range listObjectsOutput.Contents {
		if obj.Key != nil {
			objects = append(objects, *obj.Key)
		}
	}
	return objects, nil
}
