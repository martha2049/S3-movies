package main

import (
	"os"
	"path/filepath"

	"s3-movies/internal/log"
	"s3-movies/internal/s3"
	"s3-movies/internal/usecase"

	"go.uber.org/zap"

	"github.com/joho/godotenv"
)


func main() {
    // логгер
    logger, err := log.NewLogger("debug")
    if err != nil {
        logger.Fatal("cannot create logger", zap.Error(err))
        return
    }
    defer logger.Sync()

    // загружаем .env
    err = godotenv.Load()
    if err != nil {
        logger.Error("Error loading .env file", zap.Error(err))
        return
    }

	// S3 клиент
	client, err := s3.GetClient(logger)
	if err != nil {
		logger.Fatal("cannot create s3 client", zap.Error(err))
	}

	uc := usecase.NewImageUsecase(client, logger)

	srcDir := "images"

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		logger.Debug("Processing file", zap.String("file", path))

		// загрузка в S3
		_, err = uc.Upload(path)
		if err != nil {
			logger.Error("Error uploading file to S3",
				zap.String("file", path),
				zap.Error(err),
			)
			return nil
		}

		logger.Debug("Uploaded file to S3", zap.String("file", path))

		return nil
	})

	if err != nil {
		logger.Error("Error walking through images folder",
			zap.Error(err),
		)
	}

	logger.Info("All done!")
}
