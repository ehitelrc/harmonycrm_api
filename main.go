package main

import (
	"harmony_api/config"
	"harmony_api/middlewares"
	"harmony_api/routes"
	"harmony_api/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Cargar configuración
	config.LoadConfig()

	hub := ws.NewHub()
	go hub.Run()

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
