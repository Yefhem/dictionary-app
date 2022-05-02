package helpers

import (
	"io"
	"mime/multipart"
	"os"
)

func ImageUpload(fileHeader *multipart.FileHeader) (string, error) {

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create("uploads/" + fileHeader.Filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return fileHeader.Filename, nil
}
