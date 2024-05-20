package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureFileExistsAndCreate проверяет существование директории и файла, и создает их, если они не существуют.
func EnsureFileExistsAndCreate(path string) error {
	// Получаем путь к директории
	dir := filepath.Dir(path)

	// Проверяем существование директории
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Создаем директорию, если её нет
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("не удалось создать директорию: %v", err)
		}
	}

	// Проверяем существование файла
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Создаем файл, если его нет
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("не удалось создать файл: %v", err)
		}
		defer file.Close()
	}

	return nil
}
