package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	AllowedFileTypes = []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png"}
	MaxFileSize      = int64(10 * 1024 * 1024) // 10MB
)

// SaveFile saves an uploaded file to the specified directory
func SaveFile(file *multipart.FileHeader, directory string) (string, error) {
	// Validate file type
	if !ValidateFileType(file.Filename) {
		return "", fmt.Errorf("file type not allowed. Allowed types: %s", strings.Join(AllowedFileTypes, ", "))
	}

	// Validate file size
	if file.Size > MaxFileSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of 10MB")
	}

	// Create directory if it doesn't exist
	uploadDir := filepath.Join("file", directory)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return filePath, nil
}

// DeleteFile removes a file from storage
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, no error
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// ValidateFileType checks if the file extension is allowed
func ValidateFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range AllowedFileTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// GetFileExtension returns the file extension in lowercase
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// GetMimeType returns the MIME type based on file extension
func GetMimeType(filename string) string {
	ext := GetFileExtension(filename)
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}
