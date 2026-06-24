package main

import (
	"fmt"
	"harmony_api/config"
	"harmony_api/controllers"
	"harmony_api/middlewares"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/routes"
	"harmony_api/ws"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠ No se pudo cargar .env, usando variables del sistema")
	}

	r := gin.Default()

	// Cargar configuración
	config.LoadConfig()

	hub := ws.NewHub()
	go hub.Run()

	// Iniciar worker de estados
	go startStatusSyncWorker(hub)
	go startTemplateStatusSyncWorker()

	// Middleware CORS
	r.Use(middlewares.CORSMiddleware())

	routes.InitializeRoutes(r, hub)

	// Rutas WebSocket
	r.GET("/ws", ws.ServeWS(hub))
	r.Static("/static", "./uploads")

	// Iniciar el servidor en el puerto 8098
	if err := r.Run(":8098"); err != nil {
		panic("Error al iniciar el servidor: " + err.Error())
	}

}

func startStatusSyncWorker(hub *ws.Hub) {
	repo := repository.MessageRepository{}
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	fmt.Println("🚀 Iniciando worker de sincronización de estados (15s)...")

	for range ticker.C {
		fmt.Println("⏰ Tick worker estados...")
		if err := repo.ProcessUnappliedMessageStatuses(hub); err != nil {
			fmt.Printf("⚠️ Error en worker de estados: %v\n", err)
		}
	}
}

func startTemplateStatusSyncWorker() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	fmt.Println("🚀 Iniciando worker de sincronización de plantillas Meta (30m)...")

	// Run once at start
	go syncPendingTemplates()

	for range ticker.C {
		fmt.Println("⏰ Tick worker plantillas Meta...")
		syncPendingTemplates()
	}
}

func syncPendingTemplates() {
	var pendingTemplates []models.MessageTemplate
	err := config.DB.Where("approval_status = ? OR approval_status = ?", "pending", "local").Find(&pendingTemplates).Error
	if err != nil {
		fmt.Printf("⚠️ Error buscando plantillas pendientes: %v\n", err)
		return
	}

	if len(pendingTemplates) == 0 {
		return
	}

	globalWabaID, _ := repository.GetSettingTextValue("WAB_ID")

	for _, tmpl := range pendingTemplates {
		// Omit simulated templates
		if tmpl.MetaTemplateID != nil && strings.HasPrefix(*tmpl.MetaTemplateID, "simulated_id_") {
			continue
		}

		// Find access token and waba id
		var wabaID string
		var channel models.Channel
		if err := config.DB.Where("id = ?", tmpl.ChannelID).First(&channel).Error; err == nil && channel.MetaWabaID != nil && *channel.MetaWabaID != "" {
			wabaID = *channel.MetaWabaID
		} else {
			wabaID = globalWabaID
		}

		var accessToken string
		var integration models.ChannelIntegration
		if err := config.DB.Where("channel_id = ? AND is_active = ? AND access_token IS NOT NULL AND access_token != ''", tmpl.ChannelID, true).First(&integration).Error; err == nil {
			accessToken = integration.AccessToken
		}

		if accessToken == "" || wabaID == "" {
			continue
		}

		fmt.Printf("🔄 Sincronizando plantilla '%s' (ID: %d) con Meta...\n", tmpl.TemplateName, tmpl.ID)
		updates, err := controllers.SyncAndFetchMetaTemplateDetails(wabaID, accessToken, tmpl.TemplateName)
		if err != nil {
			fmt.Printf("⚠️ Error sincronizando plantilla '%s': %v\n", tmpl.TemplateName, err)
			continue
		}

		err = config.DB.Model(&tmpl).Updates(updates).Error
		if err != nil {
			fmt.Printf("⚠️ Error guardando actualización de plantilla '%s': %v\n", tmpl.TemplateName, err)
		} else {
			fmt.Printf("✅ Plantilla '%s' sincronizada con Meta. Estado: %s\n", tmpl.TemplateName, updates["approval_status"])
		}
	}
}
