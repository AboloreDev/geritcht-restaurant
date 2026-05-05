package providers

import (
	"context"
	"mime/multipart"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUploader struct {
	cld    *cloudinary.Cloudinary
	folder string
}

func NewCloudinaryUploader(cfg *config.Config) *CloudinaryUploader {
	cld, err := cloudinary.NewFromParams(
		cfg.Cloudinary.CloudinaryName,
		cfg.Cloudinary.CloudinaryAPIKey,
		cfg.Cloudinary.CloudinarySecret,
	)
	if err != nil {
		panic("Failed to create Cloudinary config" + err.Error())
	}

	return &CloudinaryUploader{
		cld:    cld,
		folder: cfg.Cloudinary.CloudinaryFolder,
	}
}

func (c *CloudinaryUploader) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ctx := context.Background()

	result, err := c.cld.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID: path,
		Folder:   c.folder,
	})
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil

}

func (c *CloudinaryUploader) DeleteFile(path string) error {
	ctx := context.Background()

	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: path,
	})

	return err
}
