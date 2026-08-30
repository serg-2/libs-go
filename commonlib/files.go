package commonlib

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// ScanDirectory рекурсивно обходит указанную папку и собирает пути к файлам по списку расширений.
// Расширения можно передавать как с точкой, так и без (например, ".txt" или "png").
func ScanDirectory(root string, extensions []string) ([]string, error) {
	var fileList []string

	// Создаем карту для быстрой проверки расширений за O(1)
	validExtensions := make(map[string]bool, len(extensions))
	for _, ext := range extensions {
		// Приводим к нижнему регистру и добавляем точку в начало, если её нет
		cleanedExt := strings.ToLower(ext)
		if !strings.HasPrefix(cleanedExt, ".") {
			cleanedExt = "." + cleanedExt
		}
		validExtensions[cleanedExt] = true
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Если это файл, проверяем его расширение
		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if validExtensions[ext] {
				fileList = append(fileList, path)
			}
		}

		return nil
	})

	return fileList, err
}
