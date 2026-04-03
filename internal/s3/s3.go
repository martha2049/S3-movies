package s3

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func GetClient() *minio.Client {
	endpoint := os.Getenv("S3_ENDPOINT")         
	accessKey := os.Getenv("S3_ACCESS_KEY")     
	secretKey := os.Getenv("S3_SECRET_KEY")      
	secure := os.Getenv("S3_SECURE") == "true"  

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		log.Fatal(err)
	}

	return client
}