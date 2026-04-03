package converter

import (
	"io/fs"
	"os"
	"path/filepath"
	"go.uber.org/zap" 
)

// конвертация из srcDir в jpeg и сохранение в dstDir
func ProcessDirectory(srcDir, dstDir string, logger *zap.Logger) error {
	// создание папки processed, если ее нет
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// чек по всем файлам
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		_, err = ConvertAndCompress(path, dstDir, logger) // Добавляем логгер
		return err
	})
}