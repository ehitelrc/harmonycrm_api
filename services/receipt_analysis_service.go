package services

import (
	"context"
	"fmt"
	"regexp"

	"harmony_api/dto"
	"harmony_api/providers"
	"harmony_api/repository"
	"harmony_api/utils"
)

type ReceiptAnalysisService struct {
	ocrService     *OCRService
	openAIProvider *providers.OpenAIProvider
	repo           *repository.ReceiptRepository
}

func NewReceiptAnalysisService(
	ocr *OCRService,
	openAI *providers.OpenAIProvider,
	repo *repository.ReceiptRepository,
) *ReceiptAnalysisService {
	return &ReceiptAnalysisService{
		ocrService:     ocr,
		openAIProvider: openAI,
		repo:           repo,
	}
}

// ────────────────────────────────────────────────────────────────
//
//	FLUJO COMPLETO: recibe base64 → OCR → extracción → normaliza
//
// ────────────────────────────────────────────────────────────────
func (s *ReceiptAnalysisService) AnalyzeFromBase64(ctx context.Context, base64Image string, caseID *uint, save bool) (*dto.ReceiptExtractionResult, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	// 1) OCR
	rawText, err := s.ocrService.ProcessBase64(base64Image)
	if err != nil {
		return nil, fmt.Errorf("error en OCR: %w", err)
	}

	// 2) Extracción semántica (OpenAI)
	result, err := s.AnalyzeFromText(ctx, rawText, nil, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ────────────────────────────────────────────────────────────────
//
//	Permite analizar texto OCR directamente (interno/reutilizable)
//
// ────────────────────────────────────────────────────────────────
func (s *ReceiptAnalysisService) AnalyzeFromText(ctx context.Context, ocrText string, caseID *uint, save bool) (*dto.ReceiptExtractionResult, error) {

	// 1) Extraer campos con OpenAI
	result, err := s.openAIProvider.ExtractReceiptData(ctx, ocrText)
	if err != nil {
		return nil, fmt.Errorf("error extrayendo datos de recibo: %w", err)
	}

	// 2) Normalizar fecha (yyyy/MM/dd)
	if result.Date != "" {
		result.Date = utils.NormalizeDate(result.Date)
	}

	// 3) Normalizar hora (HH:mm)
	if result.Time != "" {
		result.Time = utils.NormalizeTime(result.Time)
	}

	// 4) Garantizar que raw_text siempre quede lleno
	if result.RawText == "" {
		result.RawText = ocrText
	}

	// 5) Opcional: normalizar montos (luego si quieres podemos agregar NormalizeAmount)
	// result.Amount = utils.NormalizeAmount(result.Amount)
	// result.AmountSent = utils.NormalizeAmount(result.AmountSent)

	// Fix para comprobantes del BCR: Si toma un número de referencia de 8 dígitos pero existe uno de 25 terminado en esos 8 dígitos, se reemplaza.
	if result.BankName == "BCR" && len(result.ReferenceNumber) == 8 {
		re := regexp.MustCompile(`\d{17}` + regexp.QuoteMeta(result.ReferenceNumber))
		matches := re.FindStringSubmatch(ocrText)
		if len(matches) > 0 {
			result.ReferenceNumber = matches[0]
		}
	}

	// 6) Opcional:		// 4. Guardar en Base de Datos de resultados
	if save && caseID != nil {
		_, err := s.repo.SaveReceiptResult(result, *caseID, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("error guardando extracción de recibo: %w", err)
		}
	}

	return result, nil
}
