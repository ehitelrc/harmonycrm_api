package models

type OcrReportSummary struct {
	TotalCases        int64   `json:"total_cases"`
	OcrCases          int64   `json:"ocr_cases"`
	OcrPercentage     float64 `json:"ocr_percentage"`
	TotalOcrReceipts  int64   `json:"total_ocr_receipts"`
	MatchedReceipts   int64   `json:"matched_receipts"`
	UnmatchedReceipts int64   `json:"unmatched_receipts"`
}

type OcrDistribution struct {
	MessageCount int64 `json:"message_count"`
	CaseCount    int64 `json:"case_count"`
}

type OcrReportResponse struct {
	Summary      OcrReportSummary       `json:"summary"`
	Distribution []OcrDistribution      `json:"distribution"`
	Validations  []VwPaymentValidations `json:"validations"`
}
