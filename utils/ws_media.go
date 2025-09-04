package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"harmony_api/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type WSMediaMessage struct{}

func (c *WSMediaMessage) GetMediaData(url string) (string, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/json")

	auth, err := c.GetAuth()
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Authorization", *auth)

	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	var responseData models.WhatsappImageData
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", "", fmt.Errorf("error deserializing JSON: %w", err)
	}

	data, err := c.GetMediaDataFromURL(responseData.URL)
	if err != nil {
		return "", "", err
	}

	return responseData.URL, data, nil
}

func (c *WSMediaMessage) GetMediaDataFromURL(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	auth, err := c.GetAuth()
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", *auth)

	res, err := client.Do(req)
	fmt.Println("Response status:", res.Status)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(body)
	return encoded, nil
}

func (c *WSMediaMessage) GetAuth() (*string, error) {

	value := "Bearer " + "EAAXIokZBEA48BPdCvqw2Qv0e8cZAOMPQZBu4zv3lZAPN2b9gVYJZAkpu8YJH3446C8H2ssgzV9cKReGNHUyL2e30VkzsFpXk94t834GziZBgnzHTzd5dGvREyIZCvZAin9MVRKNSVZB4wP2Kcw2SUF3ITDm5NxsfqZAT7xgZAdliHO0MzpSFSKqoWeMSdJfR1WBmxNec3gUud4GUHtVtBkIAWnxQ4vsJPyiQmxee20ZAIzEp"
	return &value, nil
}

// UploadBase64AndGetURL guarda un base64 como archivo temporal y devuelve la "URL" local simulada
func UploadBase64AndGetURL(base64Content, mimeType string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return "", err
	}

	// Directorio temporal
	dir := "uploads"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Nombre de archivo con timestamp
	ext := ".bin"
	if mimeType == "image/png" {
		ext = ".png"
	} else if mimeType == "image/jpeg" {
		ext = ".jpg"
	} else if mimeType == "audio/mpeg" {
		ext = ".mp3"
	}
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(dir, fileName)

	// Guardar archivo
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}

	// Devuelve un "URL" simulada (en server real sería S3 o CDN)
	return "/static/" + fileName, nil
}
