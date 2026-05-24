package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const (
	MaxImageSize    = 10 * 1024 * 1024
	MaxImageWidth   = 1920
	MaxImageHeight  = 1080
	UploadDirectory = "uploads"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func UploadImage(file *multipart.FileHeader) (string, error) {
	if err := validateImage(file); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadPath := filepath.Join(UploadDirectory, newFilename)

	if err := os.MkdirAll(UploadDirectory, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	img = resizeImage(img)

	dst, err := os.Create(uploadPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer dst.Close()

	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(dst, img, &jpeg.Options{Quality: 85})
	case ".png":
		err = png.Encode(dst, img)
	default:
		err = jpeg.Encode(dst, img, &jpeg.Options{Quality: 85})
	}

	if err != nil {
		os.Remove(uploadPath)
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return newFilename, nil
}

func validateImage(file *multipart.FileHeader) error {
	if file.Size > MaxImageSize {
		return fmt.Errorf("file size exceeds maximum limit of 10MB")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("unsupported file type, only JPG and PNG are allowed")
	}

	return nil
}

func resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= MaxImageWidth && height <= MaxImageHeight {
		return img
	}

	return imaging.Fit(img, MaxImageWidth, MaxImageHeight, imaging.Lanczos)
}

func DeleteImage(filename string) error {
	if filename == "" {
		return nil
	}

	filePath := filepath.Join(UploadDirectory, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(filePath)
}

func ServeImage(w io.Writer, filename string) error {
	filePath := filepath.Join(UploadDirectory, filename)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	return err
}
