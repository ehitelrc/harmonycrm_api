package services

import (
	"fmt"
	"harmony_api/providers"
)

type OCRService struct {
	Provider *providers.GoogleOCR
}

func NewOCRService(provider *providers.GoogleOCR) *OCRService {
	return &OCRService{Provider: provider}
}

func (s *OCRService) ProcessBase64(base64Image string) (string, error) {
	if s.Provider == nil {
		return "", fmt.Errorf("OCR provider is not initialized (credentials file might be missing)")
	}
	return s.Provider.OCRFromBase64(base64Image)
}

func (s *OCRService) ProcessBytes(raw []byte) (string, error) {
	if s.Provider == nil {
		return "", fmt.Errorf("OCR provider is not initialized (credentials file might be missing)")
	}
	return s.Provider.OCRFromBytes(raw)
}
