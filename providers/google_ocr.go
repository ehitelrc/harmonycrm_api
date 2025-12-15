package providers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"
)

type GoogleOCR struct {
	client *vision.ImageAnnotatorClient
}

func NewGoogleOCR(credentialsPath string) (*GoogleOCR, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, fmt.Errorf("error creando cliente Vision: %v", err)
	}

	return &GoogleOCR{client: client}, nil
}

// Función que recibe base64 y devuelve texto
func (g *GoogleOCR) OCRFromBase64(base64Image string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", fmt.Errorf("error decodificando base64: %v", err)
	}

	return g.OCRFromBytes(raw)
}

// Función que recibe []byte y devuelve texto (reutilizable internamente)
func (g *GoogleOCR) OCRFromBytes(imageBytes []byte) (string, error) {
	ctx := context.Background()

	img, err := vision.NewImageFromReader(bytes.NewReader(imageBytes))

	if err != nil {
		return "", fmt.Errorf("error creando imagen desde bytes: %v", err)
	}

	annotation, err := g.client.DetectDocumentText(ctx, img, nil)
	if err != nil {
		return "", fmt.Errorf("error realizando OCR: %v", err)
	}

	if annotation == nil {
		return "", fmt.Errorf("OCR no devolvió texto")
	}

	return annotation.Text, nil
}
