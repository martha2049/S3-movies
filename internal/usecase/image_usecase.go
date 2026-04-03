package usecase

import (
	"context"
	"path/filepath"

	"s3-movies/internal/converter"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type ImageUsecase struct {
	client *minio.Client
	logger *zap.Logger
	bucket string
}

func NewImageUsecase(client *minio.Client, logger *zap.Logger) *ImageUsecase {
	return &ImageUsecase{
		client: client,
		logger: logger,
		bucket: "images",
	}
}

// публичный метод для логгера
func (u *ImageUsecase) Logger() *zap.Logger {
	return u.logger
}

// публичный метод для клиента
func (u *ImageUsecase) Client() *minio.Client {
	return u.client
}

// публичный метод для бакета
func (u *ImageUsecase) Bucket() string {
	return u.bucket
}

// обработка изображения, загрузка в S3, возврат пути к обработанному файлу
func (u *ImageUsecase) Upload(filePath string) (string, error) {
	// конвертация и сжатие изображения
	processedPath, err := converter.ConvertAndCompress(filePath, converter.ProcessedDir, u.logger)
	if err != nil {
		return "", err
	}

	// загрузка файла в S3
	_, err = u.client.FPutObject(
		context.Background(),
		u.bucket,
		filepath.Base(processedPath),
		processedPath,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return "", err
	}

	return processedPath, nil
}

// возврат списка всех объектов в бакете S3
func (u *ImageUsecase) List() ([]string, error) {
	var result []string

	objectCh := u.client.ListObjects(context.Background(), u.bucket, minio.ListObjectsOptions{Recursive: true})

	for obj := range objectCh {
		if obj.Err != nil {
			u.logger.Error("list error", zap.Error(obj.Err))
			continue
		}
		result = append(result, obj.Key)
	}

	return result, nil
}

// возврат локального пути к обработанному файлу
func (u *ImageUsecase) Download(fileName string) string {
	return filepath.Join(converter.ProcessedDir, fileName)
}