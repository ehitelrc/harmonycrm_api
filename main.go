package main

import (
	"fmt"
	"harmony_api/config"
	"harmony_api/middlewares"
	"harmony_api/repository"
	"harmony_api/routes"
	"harmony_api/ws"
	"log"
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
