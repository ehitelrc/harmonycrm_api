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

	fmt.Println("🔍 [OCR Service] Llamando a ProcessBase64 en OCR service...")
	// 1) OCR
	rawText, err := s.ocrService.ProcessBase64(base64Image)
	if err != nil {
		fmt.Printf("❌ [OCR Service] Falló ProcessBase64: %v\n", err)
		return nil, fmt.Errorf("error en OCR: %w", err)
	}

	fmt.Printf("🔍 [OCR Service] ProcessBase64 completado. Texto extraído: %d caracteres.\n", len(rawText))
	if len(rawText) > 0 {
		preview := rawText
		if len(preview) > 100 {
			preview = preview[:100]
		}
		fmt.Printf("🔍 [OCR Service] Vista previa del texto: %q\n", preview)
	}

	fmt.Println("🔍 [OCR Service] Enviando texto a OpenAI para extracción semántica...")
	// 2) Extracción semántica (OpenAI)
	result, err := s.AnalyzeFromText(ctx, rawText, nil, false)
	if err != nil {
		fmt.Printf("❌ [OCR Service] Falló la extracción con OpenAI: %v\n", err)
		return nil, err
	}

	if result != nil {
		fmt.Printf("🔍 [OCR Service] Extracción completada. Banco: %s, Referencia: %s, Monto: %f\n", result.BankName, result.ReferenceNumber, result.Amount)
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
