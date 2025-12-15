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

type WhatsAppMediaMeta struct {
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
	Sha256   string `json:"sha256"`
	FileSize int    `json:"file_size"`
}

type WSMediaMessage struct{}

// func (c *WSMediaMessage) GetMediaData(url string, channel models.ViewChannelIntegration) (string, string, error) {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	req.Header.Add("Content-Type", "application/json")

// 	auth, err := c.GetAuth(channel.AccessToken)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	req.Header.Add("Authorization", *auth)

// 	res, err := client.Do(req)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	defer res.Body.Close()

// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	var responseData models.WhatsappImageData
// 	if err := json.Unmarshal(body, &responseData); err != nil {
// 		return "", "", fmt.Errorf("error deserializing JSON: %w", err)
// 	}

// 	data, err := c.GetMediaDataFromURL(responseData.URL, channel)
// 	if err != nil {
// 		return "", "", err
// 	}

//		return responseData.URL, data, nil
//	}
func (c *WSMediaMessage) GetMediaMetadata(mediaID string, channel models.ViewChannelIntegration) (*WhatsAppMediaMeta, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v24.0/%s", mediaID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	auth, err := c.GetAuth(channel.AccessToken)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", *auth)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("media metadata error %d: %s", res.StatusCode, string(b))
	}

	var meta WhatsAppMediaMeta
	if err := json.NewDecoder(res.Body).Decode(&meta); err != nil {
		return nil, err
	}

	if meta.URL == "" {
		return nil, fmt.Errorf("media URL vacío desde WhatsApp")
	}

	return &meta, nil
}

func (c *WSMediaMessage) GetMediaDataFromURL(url string, channel models.ViewChannelIntegration) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	auth, err := c.GetAuth(channel.AccessToken)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", *auth)

	res, err := client.Do(req)

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

func (c *WSMediaMessage) GetAuth(accessToken string) (*string, error) {

	value := "Bearer " + accessToken
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
