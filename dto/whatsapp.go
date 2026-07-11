package dto

// -------------------------------------
// BASE STRUCTS REUSABLE FOR ALL FLOWS
// -------------------------------------

type WhatsAppWebhookRequest struct {
	Entry []WhatsAppEntry `json:"entry"`
}

// JSON EXACTO del webhook oficial de Meta y JSON limpio que espera tu mapper.
type WhatsAppEntry struct {
	ID      string           `json:"id"`
	Changes []WhatsAppChange `json:"changes"`
}

type WhatsAppChange struct {
	Value WhatsAppValue `json:"value"`
	Field string        `json:"field,omitempty"`
}

type WhatsAppValue struct {
	MessagingProduct string            `json:"messaging_product"`
	Metadata         WhatsAppMetadata  `json:"metadata"`
	Contacts         []WhatsAppContact `json:"contacts"`
	Messages         []WhatsAppMessage `json:"messages"`
	Statuses         []WhatsAppStatus  `json:"statuses,omitempty"`

	// Message Template Status Update fields
	Event                   string      `json:"event,omitempty"`
	MessageTemplateID       interface{} `json:"message_template_id,omitempty"` // Can be string or int64 from Meta
	MessageTemplateName     string      `json:"message_template_name,omitempty"`
	MessageTemplateLanguage string      `json:"message_template_language,omitempty"`
	Reason                  string      `json:"reason,omitempty"`
}

type WhatsAppStatus struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Timestamp    string `json:"timestamp"`
	RecipientID  string `json:"recipient_id"`
	Conversation struct {
		ID     string `json:"id"`
		Origin struct {
			Type string `json:"type"`
		} `json:"origin"`
	} `json:"conversation,omitempty"`
	Pricing struct {
		Billable     bool   `json:"billable"`
		PricingModel string `json:"pricing_model"`
		Category     string `json:"category"`
	} `json:"pricing,omitempty"`
	Errors []struct {
		Code      int    `json:"code"`
		Title     string `json:"title"`
		Message   string `json:"message"`
		ErrorData struct {
			Details string `json:"details"`
		} `json:"error_data,omitempty"`
	} `json:"errors,omitempty"`
}

type WhatsAppMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type WhatsAppContact struct {
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
	WaID string `json:"wa_id"`
}

type WhatsAppMessage struct {
	ID        string           `json:"id"`
	From      string           `json:"from"`
	Timestamp string           `json:"timestamp"`
	Type      string           `json:"type"`
	Text      WhatsAppText     `json:"text,omitempty"`
	Image     WhatsAppImage    `json:"image,omitempty"`
	Audio     WhatsAppAudio    `json:"audio,omitempty"`
	Video     WhatsAppVideo    `json:"video,omitempty"`
	Document  WhatsAppDocument `json:"document,omitempty"`
}

// -------------------------------------
// MESSAGE SUBTYPES
// -------------------------------------

type WhatsAppText struct {
	Body string `json:"body"`
}

type WhatsAppImage struct {
	Caption string `json:"caption"`
	MIME    string `json:"mime_type"`
	ID      string `json:"id"`
}

type WhatsAppAudio struct {
	ID    string `json:"id"`
	MIME  string `json:"mime_type"`
	Voice bool   `json:"voice"`
}

type WhatsAppVideo struct {
	ID      string `json:"id"`
	MIME    string `json:"mime_type"`
	Caption string `json:"caption"`
}

type WhatsAppDocument struct {
	ID       string `json:"id"`
	MIME     string `json:"mime_type"`
	Filename string `json:"filename"`
}

// -------------------------------------
// WRAPPER FROM N8N OR MANUAL REPLAY
// -------------------------------------

// ESTE DTO es para el JSON COMPLETO que entrega N8N.
// N8N envía la estructura completa, incluyendo "object", así:
// { "object": "...", "entry": [ ... ] }
//
// Solo este wrapper permite procesar mensajes pegados, pruebas Postman,
// reenvíos manuales, etc.
type N8nWhatsAppWrapper struct {
	Object string          `json:"object"`
	Entry  []WhatsAppEntry `json:"entry"`
}
