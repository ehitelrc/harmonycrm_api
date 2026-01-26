package dto

// MessengerWebhookRequest representa el payload que envía Meta (Messenger / IG)
type MessengerWebhookRequest struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`   // Page ID
		Time      int64  `json:"time"` // Epoch
		Messaging []struct {
			Sender struct {
				ID string `json:"id"` // PSID (usuario)
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"` // Page ID
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`

			Message *struct {
				Mid  string `json:"mid"`
				Text string `json:"text,omitempty"`

				Attachments []struct {
					Type    string `json:"type"` // image, audio, video, file
					Payload struct {
						URL string `json:"url,omitempty"`
					} `json:"payload"`
				} `json:"attachments,omitempty"`
			} `json:"message,omitempty"`

			Postback *struct {
				Title   string `json:"title"`
				Payload string `json:"payload"`
			} `json:"postback,omitempty"`
		} `json:"messaging"`
	} `json:"entry"`
}
