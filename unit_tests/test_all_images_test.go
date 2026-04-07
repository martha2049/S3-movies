package tests

import (
	"bytes"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/webp"
)

func TestAllImages(t *testing.T) {
	// определяем путь к файлу теста
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot get current test file path")
	}

	// images лежит в корне проекта, который один уровень выше unit_tests
	projectRoot := filepath.Dir(filepath.Dir(filename))
	imagesDir := filepath.Join(projectRoot, "images")

	entries, err := os.ReadDir(imagesDir)
	if err != nil {
		t.Fatalf("Failed to read images directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(imagesDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("%s: Failed to read file: %v", entry.Name(), err)
			continue
		}

		// декод jpeg/png
		_, format, err := image.Decode(bytes.NewReader(data))
		if err == nil {
			t.Logf("%s: Decoded successfully! Format: %s", entry.Name(), format)
			continue
		}

		// декод webp
		_, err = webp.Decode(bytes.NewReader(data))
		if err != nil {
			t.Errorf("%s: Failed to decode (jpeg/png/webp): %v", entry.Name(), err)
		} else {
			t.Logf("%s: Decoded successfully as WebP!", entry.Name())
		}
	}
}