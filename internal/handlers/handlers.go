package handlers

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"s3-movies/internal/usecase"
	"github.com/minio/minio-go/v7"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// главный endpoint
func RootHandler(c *fiber.Ctx) error {
	return c.SendString("S3 service is running")
}

// healthcheck
func HealthHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}

// загрузка файла в S3
func UploadHandler(uc *usecase.ImageUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != "POST" {
			uc.Logger().Warn("Upload attempt with wrong method", zap.String("method", c.Method()))
			return c.Status(fiber.StatusMethodNotAllowed).SendString("Method not allowed")
		}

		// загрузка файла
		file, err := c.FormFile("file")
		if err != nil {
			uc.Logger().Error("Failed to get file from form", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).SendString("Failed to upload file")
		}

		// сохранение файла
		tempFile := filepath.Join("images", file.Filename)
		if err := c.SaveFile(file, tempFile); err != nil {
			uc.Logger().Error("Failed to save file", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to upload file")
		}

		// конвертиация и сохранение
		processedPath, err := uc.Upload(tempFile) // Upload из usecase
		if err != nil {
			uc.Logger().Error("Failed to process image", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to process image")
		}

		uc.Logger().Info("File uploaded successfully", zap.String("file", filepath.Base(processedPath)))
		return c.SendString("File uploaded successfully")
	}
}

// скачивание файла из S3
func DownloadHandler(uc *usecase.ImageUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fileName := c.Query("file")
		if fileName == "" {
			uc.Logger().Warn("Download attempt missing file param")
			return c.Status(fiber.StatusBadRequest).SendString("Missing file parameter")
		}

		// скачивание из папки processed/
		localPath := filepath.Join("processed", fileName)
		if _, err := os.Stat(localPath); err != nil {
			uc.Logger().Warn("Requested file does not exist", zap.String("file", fileName))
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}

		uc.Logger().Info("Served processed file", zap.String("file", fileName))
		return c.SendFile(localPath)
	}
}

// вывод списка файлов в бакете
func ListHandler(uc *usecase.ImageUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		objectCh := uc.Client().ListObjects(context.Background(), uc.Bucket(), minio.ListObjectsOptions{Recursive: true})

		var result []string
		for obj := range objectCh {
			if obj.Err != nil {
				uc.Logger().Error("Listing object error", zap.Error(obj.Err))
				continue
			}
			result = append(result, obj.Key)
		}

		if len(result) == 0 {
			return c.SendString("No files in bucket")
		}

		return c.SendString("Files in bucket:\n" + strings.Join(result, "\n"))
	}
}