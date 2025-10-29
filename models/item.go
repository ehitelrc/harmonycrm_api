package models

type Item struct {
	ID          int     `json:"id"`
	CompanyID   int     `json:"company_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Type        string  `json:"type"` // "product" | "service"
	ItemPrice   float64 `json:"item_price"`
}

func (Item) TableName() string {
	return "items"
}
