package main

import (
	stdlog "log"
	"net/http"

	"s3-films/internal/handlers"
	"s3-films/internal/image"
	"s3-films/internal/log"
	"s3-films/internal/s3"

	"go.uber.org/zap"
)

func main() {
	// инициализация логгера
	logger, err := log.NewLogger("debug")
	if err != nil {
		stdlog.Fatalf("cannot create logger: %v", err)
	}
	defer logger.Sync()

	// S3 клиент
	client := s3.GetClient()

	// обработать исходники и сложить в processed
	err = image.ProcessDirectory("images", "processed", logger)
	if err != nil {
		logger.Fatal("failed to process images", zap.Error(err))
	}

	// роутер
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", handlers.UploadHandler(client, "images", logger))
	mux.HandleFunc("/download", handlers.DownloadHandler(client, "images", logger))
	mux.HandleFunc("/list", handlers.ListHandler(client, "images", logger))
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/", handlers.RootHandler)

	// старт сервера
	logger.Info("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Fatal("server failed", zap.Error(err)) // лог, если сервер не запускается
	}
}
