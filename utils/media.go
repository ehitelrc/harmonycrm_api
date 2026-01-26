package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SaveBase64ImageAndGetURL
// Guarda una imagen base64 en disco y devuelve una URL pública accesible por Messenger
func SaveBase64ImageAndGetURL(base64Data string, publicBaseURL string) (string, error) {

	// Quitar prefijo data:image/...;base64,
	if strings.Contains(base64Data, ",") {
		base64Data = strings.SplitN(base64Data, ",", 2)[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	// Carpeta pública
	dir := "./public/messenger"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("img_%d.jpg", time.Now().UnixNano())
	filePath := filepath.Join(dir, fileName)

	if err := os.WriteFile(filePath, decoded, 0644); err != nil {
		return "", err
	}

	// URL pública final
	publicURL := fmt.Sprintf(
		"%s/public/messenger/%s",
		strings.TrimRight(publicBaseURL, "/"),
		fileName,
	)

	return publicURL, nil
}

func SaveBase64AudioAndGetURL(
	base64Data string,
	mimeType string,
	publicBaseURL string,
) (string, error) {

	// Quitar prefijo data:audio/...;base64,
	if strings.Contains(base64Data, ",") {
		base64Data = strings.SplitN(base64Data, ",", 2)[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	// Detectar extensión
	ext := "mp3"
	if strings.Contains(mimeType, "ogg") {
		ext = "ogg"
	} else if strings.Contains(mimeType, "wav") {
		ext = "wav"
	} else if strings.Contains(mimeType, "webm") {
		ext = "webm"
	}

	dir := "./public/messenger/audio"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("audio_%d.%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(dir, fileName)

	if err := os.WriteFile(filePath, decoded, 0644); err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf(
		"%s/public/messenger/audio/%s",
		strings.TrimRight(publicBaseURL, "/"),
		fileName,
	)

	return publicURL, nil
}
