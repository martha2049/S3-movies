package s3

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

// скачивание файлов из бакета в localDir
func SyncFromS3(client *minio.Client, bucket, localDir string, logger *zap.Logger) error {
	// создание папки downloads, если нет
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// получение списка всех объектов в бакете
	objectCh := client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{Recursive: true})

	for obj := range objectCh {
		if obj.Err != nil {
			logger.Error("Error listing object", zap.Error(obj.Err))
			continue
		}

		localPath := filepath.Join(localDir, obj.Key)

		// создание поддиректорий, если есть в ключе
		if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
			logger.Error("Failed to create subdirectory", zap.String("path", localPath), zap.Error(err))
			continue
		}

		// скачивание объекта в локальный файл
		err := client.FGetObject(context.Background(), bucket, obj.Key, localPath, minio.GetObjectOptions{})
		if err != nil {
			logger.Error("Failed to download object", zap.String("object", obj.Key), zap.Error(err))
			continue
		}

		logger.Info("Downloaded object", zap.String("object", obj.Key), zap.String("to", localPath))
	}

	return nil
}
