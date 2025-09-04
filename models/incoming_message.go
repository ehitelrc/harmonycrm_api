// models/incoming_message.go
package models

type IncomingMessage struct {
	SocialNetwork string `json:"social_network" binding:"required"` // ej: fb_messenger
	RecipientID   string `json:"recipient_id" binding:"required"`   // ID del bot/app
	SenderID      string `json:"sender_id" binding:"required"`      // ID del cliente
	Timestamp     string `json:"timestamp" binding:"required"`      // Epoch timestamp en string
	TextMessage   string `json:"text_message"`                      // Texto del mensaje
	EntryID       string `json:"entry_id" binding:"required"`       // ID de entrada de Meta
	MessageID     string `json:"message_id" binding:"required"`     // ID único del mensaje
	FirstName     string `json:"first_name"`                        // (opcional)
	LastName      string `json:"last_name"`                         // (opcional)
	ProfilePic    string `json:"profile_pic"`                       // (opcional)
	SenderType    string `json:"sender_type" binding:"required"`    // Tipo de remitente (ej: client, agent)
	MessageType   string `json:"message_type" binding:"required"`   // Tipo de mensaje (ej: text, image, file, audio)
	FileURL       string `json:"file_url"`                          // URL del archivo (opcional)
	MIMEType      string `json:"mime_type"`                         // Tipo MIME del archivo (opcional)
	MediaID       string `json:"media_id"`                          // ID del medio (opcional)
	Base64Content string `json:"base64_content"`                    // Contenido en Base64 (opcional)
}

func (m *IncomingMessage) TableName() string {
	return "incoming_messages"
}
