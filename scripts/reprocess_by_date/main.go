package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"harmony_api/config"
)

type ImageRecord struct {
	CaseID    uint      `gorm:"column:case_id"`
	Base64    string    `gorm:"column:base64_content"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func main() {
	// 1. Recibir parámetros por consola
	dateStr := flag.String("date", time.Now().Format("2006-01-02"), "Fecha para consultar (YYYY-MM-DD)")
	endpoint := flag.String("endpoint", "https://harmony.ngrok.dev/api/receipts/analyze-and-save", "URL del endpoint a consumir")
	flag.Parse()

	fmt.Println("=====================================================")
	fmt.Println("🚀 Iniciando script de reprocesamiento de recibos")
	fmt.Println("📅 Fecha objetivo:", *dateStr)
	fmt.Println("🔗 Endpoint objetivo:", *endpoint)
	fmt.Println("=====================================================")

	// 2. Conectar a la base de datos de forma segura reutilizando config
	// Esto solo abre la conexión local en el script, no afecta al servicio en producción.
	config.LoadConfig()
	defer config.CloseDB()

	var records []ImageRecord

	// 3. Ejecutar la consulta SQL indicada
	query := `
		SELECT case_id, base64_content, created_at 
		FROM public.vw_cases_images_without_receipt 
		WHERE created_at::date = ? 
		ORDER BY created_at;
	`

	log.Println("🛠 Consultando la vista en la base de datos...")
	if err := config.DB.Raw(query, *dateStr).Scan(&records).Error; err != nil {
		log.Fatalf("❌ Error ejecutando la consulta: %v", err)
	}

	log.Printf("✅ Se encontraron %d registros para procesar.", len(records))

	if len(records) == 0 {
		log.Println("🛑 No hay registros que procesar para la fecha indicada. Saliendo...")
		return
	}

	// 4. Configurar el cliente HTTP (con timeout para evitar colgar el proceso)
	client := &http.Client{Timeout: 60 * time.Second}

	successCount := 0
	errorCount := 0

	// 5. Iterar, armar el JSON y hacer POST al endpoint
	for i, rec := range records {
		log.Printf("\n[%d/%d] Procesando CaseID: %d", i+1, len(records), rec.CaseID)

		payload := map[string]interface{}{
			"case_id":    rec.CaseID,
			"base64":     rec.Base64,
			"created_at": rec.CreatedAt.Format(time.RFC3339),
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Printf("  ❌ Error convirtiendo a JSON: %v", err)
			errorCount++
			continue
		}

		req, err := http.NewRequest("POST", *endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("  ❌ Error preparando la petición HTTP: %v", err)
			errorCount++
			continue
		}
		
		// Header requerido
		req.Header.Set("Content-Type", "application/json")

		// Llamada
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("  ❌ Error ejecutando POST contra el endpoint: %v", err)
			errorCount++
			continue
		}

		if resp.StatusCode == http.StatusOK {
			log.Printf("  ✅ HTTP 200 OK - Analizado y guardado con éxito")
			successCount++
		} else {
			log.Printf("  ⚠️ Advertencia - Servidor devolvió HTTP %d", resp.StatusCode)
			errorCount++
		}
		resp.Body.Close()

		// Pausa prudencial de 1 segundo para no saturar los servicios (ej: si hay IA por detrás)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n=====================================================")
	fmt.Println("🏁 Resumen del Proceso")
	fmt.Printf("Total registros identificados: %d\n", len(records))
	fmt.Printf("✅ Completados (HTTP 200): %d\n", successCount)
	fmt.Printf("❌ Con error o timeout: %d\n", errorCount)
	fmt.Println("=====================================================")
}
