package models

type TemplateListResponse struct {
	Data   []WhatsAppTemplate `json:"data"`
	Paging PagingInfo         `json:"paging"`
}

type WhatsAppTemplate struct {
	Name             string              `json:"name"`
	PreviousCategory string              `json:"previous_category,omitempty"`
	ParameterFormat  string              `json:"parameter_format,omitempty"`
	Components       []TemplateComponent `json:"components"`
	Language         string              `json:"language"`
	Status           string              `json:"status"`
	Category         string              `json:"category"`
	SubCategory      string              `json:"sub_category,omitempty"`
	ID               string              `json:"id"`
}

type TemplateComponent struct {
	Type   string `json:"type"`
	Format string `json:"format,omitempty"`
	Text   string `json:"text,omitempty"`
}

type PagingInfo struct {
	Cursors PagingCursors `json:"cursors"`
}

type PagingCursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}
