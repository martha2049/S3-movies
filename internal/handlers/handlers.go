package handlers

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"

	"s3-films/internal/image"
)

// главный endpoint
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("S3 service is running")); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// healthcheck
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// загрузка файла в S3
func UploadHandler(client *minio.Client, bucket string, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			logger.Warn("Upload attempt with wrong method", zap.String("method", r.Method))
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to upload file", http.StatusBadRequest)
			logger.Error("Failed to get file from form", zap.Error(err))
			return
		}
		defer file.Close()

		// сохранение исходника в images/
		tempFile := filepath.Join("images", header.Filename)
		out, err := os.Create(tempFile)
		if err != nil {
			http.Error(w, "Failed to upload file", http.StatusInternalServerError)
			logger.Error("Failed to create local file", zap.Error(err))
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			http.Error(w, "Failed to upload file", http.StatusInternalServerError)
			logger.Error("Failed to write file locally", zap.Error(err))
			return
		}

		// конвертация и ресайз в processed/
		processedPath, err := image.ConvertAndCompress(tempFile, "processed", logger)
		if err != nil {
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			logger.Error("Failed to convert image", zap.Error(err))
			return
		}

		// загрузка готового jpegа в S3
		if _, err := client.FPutObject(context.Background(), bucket, filepath.Base(processedPath), processedPath, minio.PutObjectOptions{}); err != nil {
			http.Error(w, "Failed to upload file", http.StatusInternalServerError)
			logger.Error("Failed to upload file to S3", zap.Error(err))
			return
		}

		if _, err := w.Write([]byte("File uploaded successfully")); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		logger.Info("File uploaded successfully", zap.String("file", filepath.Base(processedPath)))
	}
}

// скачивание файла из S3
func DownloadHandler(client *minio.Client, bucket string, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Query().Get("file")
		if fileName == "" {
			http.Error(w, "Missing file parameter", http.StatusBadRequest)
			logger.Warn("Download attempt missing file param")
			return
		}

		// скачивание из папки processed/
		localPath := filepath.Join("processed", fileName)
		if _, err := os.Stat(localPath); err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			logger.Warn("Requested file does not exist", zap.String("file", fileName))
			return
		}

		if _, err := w.Write([]byte(localPath)); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, localPath)
		logger.Info("Served processed file", zap.String("file", fileName))
	}
}

// вывод списка файлов в бакете
func ListHandler(client *minio.Client, bucket string, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		objectCh := client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{Recursive: true})

		if _, err := w.Write([]byte("Files in bucket:\n")); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for obj := range objectCh {
			if obj.Err != nil {
				logger.Error("Listing object error", zap.Error(obj.Err))
				continue
			}
			if _, err := w.Write([]byte(obj.Key + "\n")); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		logger.Info("Listed files in bucket", zap.String("bucket", bucket))
	}
}
