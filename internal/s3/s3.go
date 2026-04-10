package s3

import (
	"fmt"
	"s3-movies/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

// создание и возврат клиента MinIO/S3 
func GetClient(cfg *config.Config, logger *zap.Logger) (*minio.Client, error) {

	endpoint := cfg.S3Endpoint
	accessKey := cfg.S3AccessKey
	secretKey := cfg.S3SecretKey
	secure := cfg.S3Secure

	// проверка на пустые обязательные переменные окружения
	if endpoint == "" || accessKey == "" || secretKey == "" {
		logger.Error("Missing required S3 environment variables", zap.String("endpoint", endpoint), zap.String("accessKey", accessKey), zap.String("secretKey", secretKey))
		return nil, fmt.Errorf("missing required S3 environment variables")
	}

	// клиент MinIO
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		// логгер для записи ошибок
		logger.Error("Failed to create S3 client", zap.Error(err))
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return client, nil
}
