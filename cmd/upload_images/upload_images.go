package main

import (
	"fmt"
	"os"
	"path/filepath"

	"s3-movies/internal/converter"
	"s3-movies/internal/s3"
	"s3-movies/internal/usecase"
	loggerpkg "s3-movies/internal/log"
	"go.uber.org/zap"

	"github.com/joho/godotenv" 
)

func main() {
	// загружаем .env
	err := godotenv.Load() 
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// логгер
	logger, err := loggerpkg.NewLogger("debug")
	if err != nil {
		fmt.Printf("cannot create logger: %v\n", err)
		return
	}
	defer logger.Sync()

	// S3 клиент
	client := s3.GetClient()
	uc := usecase.NewImageUsecase(client, logger)

	srcDir := "images"
	dstDir := "processed"

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fmt.Println("Processing file:", path)

		// конвертация/сжатие
		processedPath, err := converter.ConvertAndCompress(path, dstDir, logger)
		if err != nil {
			logger.Error("Error processing file",
				zap.String("file", path),
				zap.Error(err),
			)
			return nil
		}

		// загрузка в S3
		_, err = uc.Upload(processedPath)
		if err != nil {
			logger.Error("Error uploading file to S3",
				zap.String("file", processedPath),
				zap.Error(err),
			)
			return nil
		}

		logger.Info("Uploaded file to S3",
			zap.String("file", processedPath),
		)
		return nil
	})

	if err != nil {
		logger.Error("Error walking through images folder",
			zap.Error(err),
		)
	}

	fmt.Println("All done!")
}