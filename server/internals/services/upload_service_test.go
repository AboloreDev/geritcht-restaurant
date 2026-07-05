package services

import (
	"fmt"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockUploadProvider

type MockUploadProvider struct {
	url string
	err error
}

func (m *MockUploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	return m.url, m.err
}

func (m *MockUploadProvider) DeleteFile(path string) error {
	return m.err
}

// Helper

func newUploadService(provider *MockUploadProvider) *UploadService {
	return NewUploadServices(provider)
}

// creates a fake multipart.FileHeader
func fakeFileHeader(filename string) *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: filename,
		Header:   textproto.MIMEHeader{},
		Size:     1024,
	}
}

// UploadMenuImage Tests

func TestUploadMenuImage_Success(t *testing.T) {
	service := newUploadService(&MockUploadProvider{
		url: "https://cloudinary.com/menu/1/abc.jpg",
	})

	url, err := service.UploadMenuImage(1, fakeFileHeader("image.jpg"))

	assert.NoError(t, err)
	assert.Equal(t, "https://cloudinary.com/menu/1/abc.jpg", url)
}

func TestUploadMenuImage_InvalidExtension(t *testing.T) {
	service := newUploadService(&MockUploadProvider{})

	url, err := service.UploadMenuImage(1, fakeFileHeader("malware.exe"))

	assert.Empty(t, url)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a valid image extension")
}

func TestUploadMenuImage_ValidExtensions(t *testing.T) {
	validExtensions := []string{
		"photo.jpg",
		"photo.jpeg",
		"photo.png",
		"photo.webp",
	}

	for _, filename := range validExtensions {
		t.Run(filename, func(t *testing.T) {
			service := newUploadService(&MockUploadProvider{
				url: "https://cloudinary.com/image.jpg",
			})

			url, err := service.UploadMenuImage(1, fakeFileHeader(filename))

			assert.NoError(t, err)
			assert.NotEmpty(t, url)
		})
	}
}

func TestUploadMenuImage_UploaderFails(t *testing.T) {
	service := newUploadService(&MockUploadProvider{
		err: fmt.Errorf("cloudinary error"),
	})

	url, err := service.UploadMenuImage(1, fakeFileHeader("image.jpg"))

	assert.Empty(t, url)
	assert.Error(t, err)
}

// UploadCategoryImage Tests

func TestUploadCategoryImage_Success(t *testing.T) {
	service := newUploadService(&MockUploadProvider{
		url: "https://cloudinary.com/category/abc.png",
	})

	url, err := service.UploadCategoryImage(fakeFileHeader("category.png"))

	assert.NoError(t, err)
	assert.Equal(t, "https://cloudinary.com/category/abc.png", url)
}

func TestUploadCategoryImage_InvalidExtension(t *testing.T) {
	service := newUploadService(&MockUploadProvider{})

	url, err := service.UploadCategoryImage(fakeFileHeader("file.pdf"))

	assert.Empty(t, url)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a valid image extension")
}

// DeleteFile Tests

func TestDeleteFile_Success(t *testing.T) {
	service := newUploadService(&MockUploadProvider{})

	err := service.DeleteFile(1)

	assert.NoError(t, err)
}

func TestDeleteFile_Fails(t *testing.T) {
	service := newUploadService(&MockUploadProvider{
		err: fmt.Errorf("cloudinary delete failed"),
	})

	err := service.DeleteFile(1)

	assert.Error(t, err)
}
