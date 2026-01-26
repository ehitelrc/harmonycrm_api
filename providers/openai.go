package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"harmony_api/dto"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *openai.Client
	model  string
}

func NewOpenAIProvider(model string) (*OpenAIProvider, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY no está configurada")
	}

	client := openai.NewClient(apiKey)

	if model == "" {
		model = "gpt-4o-mini"
	}

	return &OpenAIProvider{
		client: client,
		model:  model,
	}, nil
}

func (p *OpenAIProvider) ExtractReceiptData(ctx context.Context, ocrText string) (*dto.ReceiptExtractionResult, error) {
	systemPrompt := `
Eres un analizador experto de recibos bancarios de Costa Rica y Centroamérica.
Recibes el TEXTO PURO de un recibo (resultado de OCR), que puede ser de:
- SINPE móvil
- transferencias
- depósitos
- movimientos en cajero
- ventanilla
- otros formatos bancarios

Tu tarea es extraer la información estructurada.

Reglas IMPORTANTES:
- Si un campo NO aparece o NO es claro, devuélvelo como cadena vacía "" y 0 en números.
- No hagas suposiciones de banco si no aparece claramente, pero puedes inferirlo si el estilo del texto lo deja muy evidente (por ejemplo, encabezado con BAC, BNCR, BCR, Davivienda, etc.).
- Los montos (amount, amount_sent) deben ir como número decimal en formato estándar (punto como separador decimal, sin símbolos de moneda).
- Si hay múltiples montos, usa:
  - "amount" como el monto debitado de la cuenta origen.
  - "amount_sent" como el monto enviado a la cuenta destino o teléfono (si aplica).
- Devuelve exclusivamente un JSON válido, sin texto adicional ni comentarios.

Campos requeridos en el JSON:
{
  "bank_name": string,
  "transaction_type": string,        // Ej: "SINPE móvil", "transferencia", "depósito ventanilla", "cajero", etc.
  "reference_number": string,
  "date": string,                    // Tal como aparece, no cambies el formato
  "time": string,                    // Tal como aparece, no cambies el formato
  "amount": number,                  // Monto debitado
  "amount_sent": number,             // Monto enviado, si aplica
  "sender_name": string,
  "sender_phone": string,
  "receiver_name": string,
  "receiver_phone": string,
  "origin_account": string,
  "destination_account": string,
  "description": string,
  "raw_text": string,                // Copia exacta del texto OCR completo
  "warnings": []string               // Lista de advertencias, puede ir vacía
}

Reglas IMPORTANTES:
- Si un campo NO aparece o NO es claro, devuélvelo como cadena vacía "" o 0.
- Prohibido inventar datos que no existan en el OCR.
- Devuelve EXCLUSIVAMENTE JSON válido, sin comentarios adicionales.
`

	userPrompt := fmt.Sprintf("TEXTO DEL RECIBO (OCR):\n\n%s", ocrText)

	fmt.Println("ℹ️ Enviando solicitud a OpenAI para extracción de datos de recibo...")
	fmt.Println(ocrText)

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: p.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		Temperature: 0.0,
	})
	if err != nil {
		return nil, fmt.Errorf("error llamando a OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI no devolvió ninguna respuesta")
	}

	rawJSON := resp.Choices[0].Message.Content

	var result dto.ReceiptExtractionResult
	if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
		return nil, fmt.Errorf("error parseando JSON de OpenAI: %w. Respuesta fue: %s", err, rawJSON)
	}

	// Por si el modelo no rellenó el raw_text:
	if result.RawText == "" {
		result.RawText = ocrText
	}

	return &result, nil
}
