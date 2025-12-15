package services

import (
	"harmony_api/models"
	"harmony_api/repository"
)

type ReceiptStateService struct {
	repo *repository.ReceiptRepository
}

func NewReceiptStateService(repo *repository.ReceiptRepository) *ReceiptStateService {
	return &ReceiptStateService{repo: repo}
}

func (s *ReceiptStateService) ListNew() ([]models.ReceiptResult, error) {
	return s.repo.GetNewReceipts()
}

func (s *ReceiptStateService) MarkRead(id uint) error {
	return s.repo.MarkAsRead(id)
}

func (s *ReceiptStateService) MarkProcessed(id uint) error {
	return s.repo.MarkAsProcessed(id)
}

// func (s *ReceiptStateService) GenerateReceiptPDF(id uint, req dto.ProcessReceiptRequest) (string, error) {

// 	// 1️⃣ Leer template HTML
// 	templateBytes, err := os.ReadFile("templates/receipt_template.html")
// 	if err != nil {
// 		return "", fmt.Errorf("no se pudo leer el template: %w", err)
// 	}

// 	html := string(templateBytes)

// 	// 2️⃣ Reemplazar variables
// 	replacements := map[string]string{
// 		"NOMBRE":        req.Nombre,
// 		"MONTO":         fmt.Sprintf("%.2f", req.Monto),
// 		"FECHA":         req.Fecha,
// 		"REFERENCIA":    req.NumeroReferencia,
// 		"OBSERVACIONES": req.Observaciones,
// 	}

// 	for k, v := range replacements {
// 		html = strings.ReplaceAll(html, "{{"+k+"}}", v)
// 	}

// 	// 3️⃣ Rutas temporales
// 	tmpDir := "tmp/receipts"
// 	_ = os.MkdirAll(tmpDir, 0755)

// 	// pdfPath := fmt.Sprintf("%s/receipt_%d.pdf", tmpDir, id)

// 	// // 4️⃣ Generar PDF con ChromeDP
// 	// pdf, err := utils.ConvertHTMLToPDFChromedp(reportBuffer.Bytes())
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	return pdf, nil

// }
