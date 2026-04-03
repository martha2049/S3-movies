package converter

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/image/webp"
	"github.com/nfnt/resize"
	"go.uber.org/zap"
)

const MaxWidth = 800
const JPEGQuality = 80
const ProcessedDir = "processed"

// путь к исходной картинке, конвертация в JPEG, ресайз и сохранение в outputDir
func ConvertAndCompress(inputPath, outputDir string, logger *zap.Logger) (string, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// декодирование
	img, format, err := decodeImage(file, inputPath, logger)
	if err != nil {
		return "", err
	}

	// лог формата изображения
	logger.Info("Image format", zap.String("format", format))

	// ресайз
	resized := resize.Resize(MaxWidth, 0, img, resize.Lanczos3)

	// создаём папку processed/, если нет
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		logger.Error("Failed to create processed directory", zap.Error(err))
		return "", err
	}

	// подготовка выходного файла
	outputName := filepath.Base(inputPath)
	outputName = changeExtensionToJPG(outputName)
	outputPath := filepath.Join(outputDir, outputName)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// лог имени обрабатываемого файла
	logger.Info("Processing file", zap.String("file", inputPath))

	// сохранение в jpeg
	options := jpeg.Options{Quality: JPEGQuality}
	err = jpeg.Encode(outFile, resized, &options)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

// decodeImage поддерживает jpeg, png и webp
func decodeImage(r io.Reader, inputPath string, logger *zap.Logger) (image.Image, string, error) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, "", err
	}

	// декодирование стандартными методами (jpeg, png)
	img, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err == nil {
		return img, format, nil
	}

	// лог неподдерживаемого формата
	logger.Error("Failed to decode using standard formats", zap.String("file", inputPath), zap.Error(err))

	// проверка для WebP
	img, err = webp.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		// лог ошибки для webp
		logger.Error("Failed to decode webp", zap.String("file", inputPath), zap.Error(err))
		return nil, "", err
	}

	return img, "webp", nil
}

// добавление расширения .jpg
func changeExtensionToJPG(filename string) string {
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	return name + ".jpg"
}