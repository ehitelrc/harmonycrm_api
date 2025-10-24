package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"harmony_api/models"
	"harmony_api/ws"
	"net/http"
	"strconv"
	"sync"
)

//const apiBaseURL = "https://graph.facebook.com/v21.0"

// Structs para el payload de WhatsApp
type WhatsAppMessage struct {
	MessagingProduct string        `json:"messaging_product"`
	To               string        `json:"to"`
	Type             string        `json:"type"`
	Template         TemplateBlock `json:"template"`
}

type TemplateBlock struct {
	Name     string        `json:"name"`
	Language LanguageBlock `json:"language"`
}

type LanguageBlock struct {
	Code string `json:"code"`
}

// SendTemplateMessage envía un template a un solo número
func SendTemplateMessage(apiBaseURL, phoneNumberID, accessToken, templateName, languageCode, to string) error {
	msg := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: TemplateBlock{
			Name: templateName,
			Language: LanguageBlock{
				Code: languageCode,
			},
		},
	}

	body, _ := json.Marshal(msg)

	// body to string json for logging
	fmt.Println("Payload JSON:", string(body))

	//url := fmt.Sprintf("%s/%s/messages", apiBaseURL, phoneNumberID)
	url := apiBaseURL + "/" + phoneNumberID + "/messages"

	fmt.Println("Enviando a URL:", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("error en respuesta: %s", resp.Status)
	}

	fmt.Printf("📩 Mensaje enviado a %s | Status: %s\n", to, resp.Status)
	return nil
}

// SendTemplateToMany envía un template a múltiples números en paralelo
func SendTemplateToManyx(apiBaseURL, phoneNumberID, accessToken, templateName, languageCode string, numbers []string) {
	var wg sync.WaitGroup

	for _, number := range numbers {
		wg.Add(1)
		go func(num string) {
			defer wg.Done()
			if err := SendTemplateMessage(apiBaseURL, phoneNumberID, accessToken, templateName, languageCode, num); err != nil {
				fmt.Printf("❌ Error enviando a %s: %v\n", num, err)
			}
		}(number)
	}

	wg.Wait()
	fmt.Println("✅ Todos los mensajes procesados.")
}

func SendTemplateToMany(
	apiBaseURL, phoneNumberID, accessToken, templateName, languageCode string,
	targets []models.TemplateRecipient,
	hub *ws.Hub,

) {
	var wg sync.WaitGroup

	for _, t := range targets {
		wg.Add(1)
		go func(target models.TemplateRecipient) {
			defer wg.Done()

			// 1️⃣ Enviar template al número
			if err := SendTemplateMessage(apiBaseURL, phoneNumberID, accessToken, templateName, languageCode, target.Number); err != nil {
				fmt.Printf("❌ Error enviando a %s: %v\n", target.Number, err)
				return
			}
			fmt.Printf("✅ Template enviado a %s\n", target.Number)

			// 2️⃣ Si hay case_id, emitir notificación por WebSocket
			if target.CaseID != nil && *target.CaseID > 0 && hub != nil {
				payload, _ := json.Marshal(map[string]interface{}{
					"type":    "new_message",
					"case_id": *target.CaseID,
					"data": map[string]interface{}{
						"sender_type":  "client",
						"message_type": "text",
						"text_content": "Apertura mediante template",
					},
				})

				channel := "case:" + strconv.FormatInt(*target.CaseID, 10)
				hub.BroadcastJSON(channel, payload)
				fmt.Printf("📢 WS emitido a canal %s\n", channel)
			}
		}(t)
	}

	wg.Wait()
	fmt.Println("✅ Todos los templates procesados.")
}
