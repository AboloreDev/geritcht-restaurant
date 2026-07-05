package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/google/uuid"
)

type UploadService struct {
	uploader interfaces.UploadProvider
}

func NewUploadServices(uploader interfaces.UploadProvider) *UploadService {
	return &UploadService{
		uploader: uploader,
	}
}

func (s *UploadService) UploadMenuImage(menuID uint, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !utils.IsValidExtensions(ext) {
		return "", fmt.Errorf("This is not a valid image extension %s", ext)
	}

	newFileName := uuid.New().String() + ext

	path := fmt.Sprintf("menu/%d/%s", menuID, newFileName)

	return s.uploader.UploadFile(file, path)

}

func (s *UploadService) UploadCategoryImage(file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !utils.IsValidExtensions(ext) {
		return "", fmt.Errorf("This is not a valid image extension %s", ext)
	}

	newFileName := uuid.New().String() + ext

	path := fmt.Sprintf("category/%s", newFileName)

	return s.uploader.UploadFile(file, path)
}

func (s *UploadService) DeleteFile(menuID uint) error {
	path := fmt.Sprintf("menu/%d", menuID)

	err := s.uploader.DeleteFile(path)
	if err != nil {
		return err
	}

	return nil
}
