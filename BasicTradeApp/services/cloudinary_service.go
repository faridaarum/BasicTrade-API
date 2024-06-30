package services

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var allowedFormats = map[string]bool{
	"image/jpeg":    true,
	"image/jpg":     true,
	"image/png":     true,
	"image/svg+xml": true,
}

const maxFileSize = 5 * 1024 * 1024 // 5 MB

func UploadFileToCloudinary(file multipart.File, fileHeader *multipart.FileHeader, fileName string) (string, error) {
	if fileHeader.Size > maxFileSize {
		return "", errors.New("file size exceeds the 5MB limit")
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedFormats[contentType] {
		return "", errors.New("file format is not allowed")
	}

	cld, err := cloudinary.NewFromParams("your_cloud_name", "your_api_key", "your_api_secret")
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: fileName,
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
