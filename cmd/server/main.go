package main

import (
	"context"
	stdlog "log"
	"time"

	"s3-movies/internal/handlers"
	"s3-movies/internal/log"
	"s3-movies/internal/s3"
	"s3-movies/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const MaxUploadSize = 50 * 1024 * 1024

func main() {
	logger, err := log.NewLogger("debug")
	if err != nil {
		stdlog.Fatalf("cannot create logger: %v", err)
	}
	defer logger.Sync()

	client, err := s3.GetClient(logger)
	if err != nil {
		logger.Error("cannot create S3 client, continuing without S3", zap.Error(err))
		client = nil 
	}

	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := s3.SyncFromS3(ctx, client, "images", "downloads", logger); err != nil {
			logger.Error("failed to sync S3", zap.Error(err))
		}
	}

	// usecase
	uc := usecase.NewImageUsecase(client, logger)

	// создаём Fiber с лимитом на body
	app := fiber.New(fiber.Config{
		BodyLimit: MaxUploadSize,
	})

	// раздаём статические файлы UI
	app.Static("/", "./internal/ui")

	// роуты
	app.Get("/", handlers.RootHandler)
	app.Get("/health", handlers.HealthHandler)

	// обновляем хэндлеры
	app.Post("/upload", handlers.UploadHandler(uc))
	app.Get("/download", handlers.DownloadHandler(uc))
	app.Get("/list", handlers.ListHandler(uc))

	logger.Info("Server started on :8080")

	if err := app.Listen(":8080"); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
