package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func FileName(filePath string) string {
	return filepath.Base(filePath)
}

// Ensures the file at the specified path is an image and returns its resolued
// absolute path.
func ValidateImageFile(filePath string) (string, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}
	if !isFile(absFilePath) {
		return "", fmt.Errorf("File doesn't exists: %s", absFilePath)
	}
	// Validate extension is valid.
	fileExt := strings.ToLower(filepath.Ext(absFilePath))
	switch fileExt {
	case ".png", ".jpg", ".jpeg":
		return absFilePath, nil
	default:
		return "", fmt.Errorf("Invalid file: %s", absFilePath)
	}
}
