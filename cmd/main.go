package main

import (
    stdlog "log"                
    "s3-movies/internal/log"     
    "s3-movies/internal/handlers"
    "s3-movies/internal/s3"
    "s3-movies/internal/usecase"

    "github.com/gofiber/fiber/v2"
)
const MaxUploadSize = 50 * 1024 * 1024

func main() {
	logger, err := log.NewLogger("debug")
	if err != nil {
		stdlog.Fatalf("cannot create logger: %v", err)
	}
	defer logger.Sync()

	client := s3.GetClient()

	// usecase
	uc := usecase.NewImageUsecase(client, logger)

    // создаём Fiber с лимитом на body
    app := fiber.New(fiber.Config{
        BodyLimit: MaxUploadSize,
    })

	// роуты
	app.Get("/", handlers.RootHandler)
	app.Get("/health", handlers.HealthHandler)

	// обновляем хэндлеры
	app.Post("/upload", handlers.UploadHandler(uc))
	app.Get("/download", handlers.DownloadHandler(uc))
	app.Get("/list", handlers.ListHandler(uc))

	logger.Info("Server started on :8080")
	stdlog.Fatal(app.Listen(":8080"))
}