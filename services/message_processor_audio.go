package services

import (
	"fmt"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
)

// ---------------------------------------------
// PROCESAMIENTO DE AUDIO (ASYNC, FIEL A TU CÓDIGO)
// ---------------------------------------------
func (p *MessageProcessor) processAudio(
	input models.IncomingMessage,
	newMessage *models.Message,
) {

	fmt.Println("🎧 Procesando audio async | message_id:", newMessage.ID)

	// 1️⃣ Obtener integración del canal
	channelRepository := repository.ChannelRepository{}

	channel, err := channelRepository.GetChannelIntegrationByAppIdentifier(input.RecipientID)
	if err != nil || channel == nil {
		fmt.Println("❌ Error obteniendo integración del canal:", err)
		return
	}

	wmUtils := utils.WSMediaMessage{}

	// 2️⃣ Obtener METADATA del audio (NO el binario)
	meta, err := wmUtils.GetMediaMetadata(input.MediaID, *channel)
	if err != nil {
		fmt.Println("❌ Error obteniendo metadata del audio:", err)
		return
	}

	if meta.URL == "" {
		fmt.Println("❌ Metadata inválida: URL vacía")
		return
	}

	// 3️⃣ Descargar BINARIO real desde la URL
	resourceData, err := wmUtils.GetMediaDataFromURL(meta.URL, *channel)
	if err != nil {
		fmt.Println("❌ Error descargando binario del audio:", err)
		return
	}

	// 4️⃣ Construir Base64 completo (MISMO FORMATO QUE YA USAS)
	completeData := "data:" + meta.MimeType + ";base64," + resourceData
	input.Base64Content = completeData
	input.MIMEType = meta.MimeType

	fmt.Println("✅ Audio descargado y normalizado | MIME:", meta.MimeType)

	// ⚠️ IMPORTANTE
	// El mensaje YA fue guardado por ProcessIncomingMessage
	// Aquí NO se vuelve a guardar ni se emite WS

	// 5️⃣ Punto de extensión FUTURO (tal como lo tenías)
	// - Conversión WebM → OGG
	// - Speech-to-text
	// - Clasificación

	msgRepo := repository.MessageRepository{}
	if err := msgRepo.UpdateMediaContent(
		newMessage.ID,
		input.Base64Content,
		input.MIMEType,
	); err != nil {
		fmt.Println("❌ Error guardando audio en DB:", err)
		return
	}

	fmt.Println("✅ Audio persistido correctamente | message_id:", newMessage.ID)
}
