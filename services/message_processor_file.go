package services

import (
	"fmt"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
)

func (p *MessageProcessor) processFile(
	input models.IncomingMessage,
	newMessage *models.Message,
) {

	// 1️⃣ Obtener integración del canal
	channelRepository := repository.ChannelRepository{}

	channel, err := channelRepository.GetChannelIntegrationByAppIdentifier(input.RecipientID)
	if err != nil || channel == nil {
		fmt.Println("❌ Error obteniendo integración del canal (file):", err)
		return
	}

	wmUtils := utils.WSMediaMessage{}

	// 2️⃣ OBTENER METADATA DEL MEDIA
	meta, err := wmUtils.GetMediaMetadata(input.MediaID, *channel)
	if err != nil {
		fmt.Println("❌ Error obteniendo metadata del archivo:", err)
		return
	}

	if meta.URL == "" {
		fmt.Println("❌ Metadata inválida: URL vacía")
		return
	}

	// 3️⃣ DESCARGAR BINARIO REAL
	resourceData, err := wmUtils.GetMediaDataFromURL(meta.URL, *channel)
	if err != nil {
		fmt.Println("❌ Error descargando binario del archivo:", err)
		return
	}

	// 4️⃣ Construir Base64 completo
	// Usamos el mime type que viene de meta o del input si meta falla (aunque meta es más confiable)
	mime := meta.MimeType
	if mime == "" {
		mime = input.MIMEType
	}

	completeData := "data:" + mime + ";base64," + resourceData

	// Actualizamos el input para consistencia, aunque lo importante es lo que va a la BD
	input.Base64Content = completeData

	// 5️⃣ ACTUALIZAR BD
	msgRepo := repository.MessageRepository{}

	if err := msgRepo.UpdateMediaContent(
		newMessage.ID,
		completeData,
		mime,
	); err != nil {
		fmt.Println("❌ Error actualizando base64 de archivo:", err)
		return
	}

	fmt.Println("✅ Archivo procesado y guardado:", newMessage.ID)
}
