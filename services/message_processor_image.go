package services

import (
	"fmt"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"strings"
)

func (p *MessageProcessor) processImage(
	input models.IncomingMessage,
	newMessage *models.Message,
) {

	// 1️⃣ Obtener integración del canal
	channelRepository := repository.ChannelRepository{}

	channel, err := channelRepository.GetChannelIntegrationByAppIdentifier(input.RecipientID)
	if err != nil || channel == nil {
		fmt.Println("❌ Error obteniendo integración del canal:", err)
		return
	}

	wmUtils := utils.WSMediaMessage{}

	// 2️⃣ OBTENER METADATA DEL MEDIA (NO EL BINARIO)
	meta, err := wmUtils.GetMediaMetadata(input.MediaID, *channel)
	if err != nil {
		fmt.Println("❌ Error obteniendo metadata del medio:", err)
		return
	}

	if meta.URL == "" {
		fmt.Println("❌ Metadata inválida: URL vacía")
		return
	}

	// 3️⃣ DESCARGAR BINARIO REAL
	resourceData, err := wmUtils.GetMediaDataFromURL(meta.URL, *channel)
	if err != nil {
		fmt.Println("❌ Error descargando binario del medio:", err)
		return
	}

	// 4️⃣ Construir Base64 completo
	completeData := "data:" + meta.MimeType + ";base64," + resourceData
	input.Base64Content = completeData

	// ⚠️ IMPORTANTE
	// El mensaje YA fue creado en ProcessIncomingMessage
	// aquí NO se vuelve a guardar

	msgRepo := repository.MessageRepository{}

	if err := msgRepo.UpdateMediaContent(
		newMessage.ID,
		input.Base64Content,
		input.MIMEType,
	); err != nil {
		fmt.Println("❌ Error actualizando base64 de imagen:", err)
		return
	}

	newMessage.Base64Content = input.Base64Content
	newMessage.MIMEType = input.MIMEType
	p.broadcast(newMessage)

	// 5️⃣ OCR ASÍNCRONO (idéntico a tu lógica actual)
	if !channel.AnalyzeIncomingImages {
		fmt.Println("ℹ️ Análisis de imágenes entrantes deshabilitado para esta integración.")
		return
	}

	go func(input models.IncomingMessage, newMessage *models.Message) {
		caseID := uint(newMessage.CaseID)
		fmt.Printf("🔍 [OCR Trace] Iniciando análisis para mensaje ID %d (Case ID %d)\n", newMessage.ID, caseID)

		// Base64 sin prefijo
		b64 := input.Base64Content
		if idx := strings.Index(b64, ","); idx != -1 {
			b64 = b64[idx+1:]
		}

		fmt.Printf("🔍 [OCR Trace] Llamando a AnalyzeFromBase64 con base64 de tamaño %d bytes\n", len(b64))
		result, err := p.receiptSvc.AnalyzeFromBase64(
			nil, // no gin.Context aquí
			b64,
			&caseID,
			true,
		)
		if err != nil {
			fmt.Printf("❌ [OCR Trace] Error analizando recibo para mensaje %d: %v\n", newMessage.ID, err)
			return
		}

		if result == nil {
			fmt.Printf("ℹ️ [OCR Trace] La imagen del mensaje %d no es un recibo.\n", newMessage.ID)
			return
		}

		fmt.Printf("🔍 [OCR Trace] Analizado con éxito. Banco: %s, Ref: %s, Monto: %f. Guardando en BD...\n", result.BankName, result.ReferenceNumber, result.Amount)
		receiptRepo := repository.NewReceiptRepository()

		messageID := uint64(newMessage.ID)
		record, err := receiptRepo.SaveReceiptResult(result, caseID, &messageID, nil)
		if err != nil {
			fmt.Printf("❌ [OCR Trace] Error guardando recibo en la BD para mensaje %d: %v\n", newMessage.ID, err)
			return
		}

		fmt.Printf("✅ [OCR Trace] Recibo guardado con ID: %d para mensaje %d\n", record.ID, newMessage.ID)

	}(input, newMessage)
}
