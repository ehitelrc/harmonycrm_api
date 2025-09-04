package ws

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// IMPORTANTE: Endurece esto para producción (valida orígenes permitidos)
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Si prefieres inyectar el hub desde fuera:
func ServeWS(h *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ej: /ws?case_id=123  -> canal = "case:123"
		caseID := c.Query("case_id")
		agentID := c.Query("agent_id")

		fmt.Println("Conexión WS entrante - case_id:", caseID, "agent_id:", agentID)

		channel := ""
		if caseID != "" {
			channel = "case:" + caseID
		} else if agentID != "" {
			channel = "agent:" + agentID
		} else {
			c.Status(http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := &Client{
			Conn:    conn,
			Channel: channel,
			Send:    make(chan []byte, 256),
		}

		h.register <- client

		// Writer goroutine
		go func() {
			for msg := range client.Send {
				if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					break
				}
			}
		}()

		// Reader loop (opcional si quieres recibir eventos del cliente)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}

		h.unregister <- client
	}
}
