package models

type WhatsappImageData struct {
	URL              string `json:"url"`
	MimeType         string `json:"mime_type"`
	SHA256           string `json:"sha256"`
	FileSize         int    `json:"file_size"`
	ID               string `json:"id"`
	MessagingProduct string `json:"messaging_product"`
}
