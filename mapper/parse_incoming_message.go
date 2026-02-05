package mapper

import (
	"harmony_api/dto"
	"harmony_api/models"
)

func ParseWhatsAppToIncoming(req dto.WhatsAppWebhookRequest) *models.IncomingMessage {

	entry := req.Entry[0]
	change := entry.Changes[0]
	value := change.Value
	msg := value.Messages[0]

	var (
		messageType string
		textMessage string
		fileURL     string
		mimeType    string
		mediaID     string
	)

	switch msg.Type {
	case "text":
		messageType = "text"
		textMessage = msg.Text.Body

	case "image":
		messageType = "image"
		textMessage = msg.Image.Caption // caption → text_message
		mimeType = msg.Image.MIME
		mediaID = msg.Image.ID
		// fileURL se obtiene después llamando a /media/ID

	case "audio":
		messageType = "audio"
		mimeType = msg.Audio.MIME
		mediaID = msg.Audio.ID

	case "document":
		messageType = "file"                // Internal type 'file'
		textMessage = msg.Document.Filename // filename → text_message
		mimeType = msg.Document.MIME
		mediaID = msg.Document.ID
	}

	incoming := models.IncomingMessage{
		SocialNetwork: "whatsapp",
		RecipientID:   value.Metadata.PhoneNumberID,
		SenderID:      msg.From,
		Timestamp:     msg.Timestamp,
		TextMessage:   textMessage,
		EntryID:       entry.ID,
		MessageID:     msg.ID,
		FirstName:     "",
		LastName:      "",
		ProfilePic:    "",
		SenderType:    "client",
		MessageType:   messageType,
		FileURL:       fileURL,
		MIMEType:      mimeType,
		MediaID:       mediaID,
		Base64Content: "",
	}

	return &incoming
}
