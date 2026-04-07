package s3

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// создание и возврат клиента MinIO/S3.
func GetClient(logger *zap.Logger) (*minio.Client, error) {
	// астройка Viper для работы с переменными окружения
	viper.SetEnvPrefix("S3")   
	viper.AutomaticEnv()        

	// чтение переменных окружения 
	endpoint := viper.GetString("ENDPOINT")
	accessKey := viper.GetString("ACCESS_KEY")
	secretKey := viper.GetString("SECRET_KEY")
	secure := viper.GetBool("SECURE")

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