package mapper

import (
	"strconv"

	"harmony_api/dto"
	"harmony_api/models"
)

func ParseMessengerWebhook(payload dto.MessengerWebhookRequest) []models.IncomingMessage {
	var messages []models.IncomingMessage

	for _, entry := range payload.Entry {
		for _, event := range entry.Messaging {

			if event.Message == nil {
				continue
			}

			senderID := event.Sender.ID
			recipientID := event.Recipient.ID
			messageID := event.Message.Mid
			timestamp := strconv.FormatInt(event.Timestamp, 10)

			// -------- TEXTO --------
			if event.Message.Text != "" {
				messages = append(messages, models.IncomingMessage{
					SocialNetwork: "fb_messenger",
					RecipientID:   recipientID,
					SenderID:      senderID,
					Timestamp:     timestamp,
					TextMessage:   event.Message.Text,
					EntryID:       entry.ID,
					MessageID:     messageID,
					SenderType:    "client",
					MessageType:   "text",
				})
			}

			// -------- ATTACHMENTS --------
			for _, att := range event.Message.Attachments {
				msgType := att.Type // image | audio | video | file

				messages = append(messages, models.IncomingMessage{
					SocialNetwork: "fb_messenger",
					RecipientID:   recipientID,
					SenderID:      senderID,
					Timestamp:     timestamp,
					EntryID:       entry.ID,
					MessageID:     messageID,
					SenderType:    "client",
					MessageType:   msgType,
					FileURL:       att.Payload.URL,
				})
			}
		}
	}

	return messages
}
