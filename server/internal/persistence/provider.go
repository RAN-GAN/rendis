package persistence

import (
	"errors"
	"os"
)

type Provider interface {
	Load(name string) ([]byte, error)
	Save(name string, data []byte) error
	Delete(name string) error
	List() ([]string, error)
}

func NewProvider() (Provider, error) {
	provider := os.Getenv("SNAPSHOT_PROVIDER")
	if provider == "" {
		return nil, errors.New("No Persistance provided")
	}
	switch provider {
	case "S3":
		cfg := S3Config{
			Endpoint:  os.Getenv("S3_ENDPOINT"),
			Region:    os.Getenv("S3_REGION"),
			Bucket:    os.Getenv("S3_BUCKET"),
			AccessKey: os.Getenv("S3_ACCESS_KEY_ID"),
			SecretKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
		}

		s3Client, err := NewS3Provider(cfg)
		if err != nil {
			return nil, err
		}
		return s3Client, nil
	default:
		return nil, errors.New("Failed creating provider")
	}
}
