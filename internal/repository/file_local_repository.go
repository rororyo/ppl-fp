package repository

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalFileRepository struct {
	BasePath string
	BaseURL  string
}

func NewLocalFileRepository(basePath, baseURL string) *LocalFileRepository {
	return &LocalFileRepository{
		BasePath: basePath,
		BaseURL:  baseURL,
	}
}

func (r *LocalFileRepository) UploadFile(file multipart.File, fileName string, contentType string) (string, error) {
	// Ensure the base path exists
	fullPath := filepath.Join(r.BasePath, fileName)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directories: %w", err)
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filepath.ToSlash(filepath.Join(r.BaseURL, fileName)), nil
}

func (r *LocalFileRepository) DeleteFile(fileURL string) error {
	filePath := filepath.Join(r.BasePath, filepath.Base(fileURL))
	return os.Remove(filePath)
}
