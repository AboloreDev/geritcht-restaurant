package providers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalUploadProvider struct {
	basePath string
}

func NewLocalUploadProvider(basePath string) *LocalUploadProvider {
	return &LocalUploadProvider{basePath: basePath}
}

func (u *LocalUploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	fullPath := filepath.Join(u.basePath, path)
	fullFilePath := filepath.Dir(fullPath)

	err := os.Mkdir(fullFilePath, 0755)
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	destination, err := os.Create(fullFilePath)
	if err != nil {
		return "", nil
	}
	defer destination.Close()

	_, err = destination.ReadFrom(src)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/%s", path), nil

}

func (u *LocalUploadProvider) DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}
