package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/veetmoradiya3628/go-shop/internal/interfaces"
)

type UploadService struct {
	Provider interfaces.UploadProvider
}

func NewUploadService(provider interfaces.UploadProvider) *UploadService {
	return &UploadService{Provider: provider}
}

func (s *UploadService) UploadProductImage(productID uint, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageExtension(ext) {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}
	path := fmt.Sprintf("products/%d/%s", productID, file.Filename)
	return s.Provider.UploadFile(file, path)
}

func isValidImageExtension(ext string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}
