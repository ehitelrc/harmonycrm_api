package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type OcrReportRepository struct{}

func NewOcrReportRepository() *OcrReportRepository {
	return &OcrReportRepository{}
}

func (r *OcrReportRepository) GetOcrReport(companyID int64, startDate, endDate string) (*models.OcrReportResponse, error) {
	var summary models.OcrReportSummary
	var distribution []models.OcrDistribution
	var validations []models.VwPaymentValidations

	// Date filters initialization
	dateFilterCases := ""
	dateFilterReceipts := ""
	dateFilterReceiptsLeftJoin := ""
	var argsCases []interface{}
	var argsReceipts []interface{}
	var argsReceiptsLeftJoin []interface{}

	argsCases = append(argsCases, companyID)
	argsReceipts = append(argsReceipts, companyID)
	argsReceiptsLeftJoin = append(argsReceiptsLeftJoin, companyID)

	if startDate != "" && endDate != "" {
		dateFilterCases = " AND created_at >= ? AND created_at <= ?"
		argsCases = append(argsCases, startDate+" 00:00:00", endDate+" 23:59:59.999999")

		dateFilterReceipts = " AND rr.created_at >= ? AND rr.created_at <= ?"
		argsReceipts = append(argsReceipts, startDate+" 00:00:00", endDate+" 23:59:59.999999")

		dateFilterReceiptsLeftJoin = " AND rr.created_at >= ? AND rr.created_at <= ?"
		argsReceiptsLeftJoin = append(argsReceiptsLeftJoin, startDate+" 00:00:00", endDate+" 23:59:59.999999")
	} else if startDate != "" {
		dateFilterCases = " AND created_at >= ?"
		argsCases = append(argsCases, startDate+" 00:00:00")

		dateFilterReceipts = " AND rr.created_at >= ?"
		argsReceipts = append(argsReceipts, startDate+" 00:00:00")

		dateFilterReceiptsLeftJoin = " AND rr.created_at >= ?"
		argsReceiptsLeftJoin = append(argsReceiptsLeftJoin, startDate+" 00:00:00")
	}

	// 1. Total Cases count in selected range
	totalCasesSQL := "SELECT COUNT(*) FROM cases WHERE company_id = ?" + dateFilterCases
	if err := config.DB.Raw(totalCasesSQL, argsCases...).Scan(&summary.TotalCases).Error; err != nil {
		return nil, err
	}

	// 2. OCR Cases count (Cases with receipt results) in selected range
	ocrCasesSQL := `
		SELECT COUNT(DISTINCT rr.case_id) 
		FROM receipt_results rr 
		JOIN cases c ON c.id = rr.case_id 
		WHERE c.company_id = ?` + dateFilterReceipts
	if err := config.DB.Raw(ocrCasesSQL, argsReceipts...).Scan(&summary.OcrCases).Error; err != nil {
		return nil, err
	}

	// Calculate percentage
	if summary.TotalCases > 0 {
		summary.OcrPercentage = float64(summary.OcrCases) * 100.0 / float64(summary.TotalCases)
	}

	// 3. Receipt match statistics (Confirmed in ERP vs Unconfirmed)
	receiptMatchSQL := `
		SELECT 
			COUNT(rr.id) as total_ocr_receipts,
			COUNT(rr.id) FILTER (WHERE epc.id_record IS NOT NULL) as matched_receipts,
			COUNT(rr.id) FILTER (WHERE epc.id_record IS NULL) as unmatched_receipts
		FROM receipt_results rr
		JOIN cases c ON c.id = rr.case_id
		LEFT JOIN erp_payment_confirmation epc ON epc.reference_number = rr.reference_number
		WHERE c.company_id = ?` + dateFilterReceipts
	if err := config.DB.Raw(receiptMatchSQL, argsReceipts...).Scan(&summary).Error; err != nil {
		return nil, err
	}

	// 4. Distribution of OCR messages per case
	distSQL := `
		SELECT message_count, COUNT(*) as case_count
		FROM (
			SELECT rr.case_id, COUNT(*) as message_count
			FROM receipt_results rr
			JOIN cases c ON c.id = rr.case_id
			WHERE c.company_id = ?` + dateFilterReceipts + `
			GROUP BY rr.case_id
		) sub
		GROUP BY message_count
		ORDER BY message_count ASC`
	if err := config.DB.Raw(distSQL, argsReceipts...).Scan(&distribution).Error; err != nil {
		return nil, err
	}

	// 5. Validations (using view vw_payment_validations)
	validationsSQL := `
		SELECT * 
		FROM vw_payment_validations 
		WHERE case_id IN (SELECT id FROM cases WHERE company_id = ?)` + dateFilterReceiptsLeftJoin + `
		ORDER BY ocr_date DESC`
	if err := config.DB.Raw(validationsSQL, argsReceiptsLeftJoin...).Scan(&validations).Error; err != nil {
		return nil, err
	}

	return &models.OcrReportResponse{
		Summary:      summary,
		Distribution: distribution,
		Validations:  validations,
	}, nil
}
