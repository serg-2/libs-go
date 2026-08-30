package imagelib

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	golibraw "github.com/seppedelanghe/go-libraw"
)

// DecodeImageFile - decode JPG, PNG, GIF, CR2
func DecodeImageFile(filePath string) (image.Image, error) {
	// Если это RAW-файл Canon, используем go-libraw
	if strings.ToLower(filepath.Ext(filePath)) == ".cr2" {
		img, err := GetImageFromCR2(filePath)
		if err != nil {
			return nil, fmt.Errorf("go-libraw: %w", err)
		}
		return img, nil
	}

	// Для всех стандартных форматов (PNG, JPG, GIF)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("image.Decode: %w", err)
	}

	return img, nil
}

func GetImageFromCR2(srcPath string) (image.Image, error) {
	// 1. Инициализируем процессор с настройками по умолчанию
	processor := golibraw.NewProcessor(golibraw.NewProcessorOptions())

	// 2. Вызываем корректный метод для быстрого извлечения превью
	// Метод возвращает ([]byte, PreviewInfo, error)
	previewData, _, err := processor.ExtractLargestPreview(srcPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка извлечения превью из %s: %w", srcPath, err)
	}

	// 3. Декодируем байты встроенного JPEG изображения в интерфейс image.Image
	img, _, err := image.Decode(bytes.NewReader(previewData))
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования JPEG данных: %w", err)
	}

	return img, nil
}
