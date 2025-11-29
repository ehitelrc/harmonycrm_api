package dto

type WhatsAppWebhookRequestt struct {
	Entry []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`

				Metadata struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`

				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`

				Messages []struct {
					ID        string `json:"id"`
					From      string `json:"from"`
					Timestamp string `json:"timestamp"`
					Type      string `json:"type"`

					// TEXT
					Text struct {
						Body string `json:"body"`
					} `json:"text,omitempty"`

					// IMAGE
					Image struct {
						Caption string `json:"caption"`
						MIME    string `json:"mime_type"`
						ID      string `json:"id"`
					} `json:"image,omitempty"`

					// AUDIO
					Audio struct {
						ID    string `json:"id"`
						MIME  string `json:"mime_type"`
						Voice bool   `json:"voice"`
					} `json:"audio,omitempty"`
				} `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}
