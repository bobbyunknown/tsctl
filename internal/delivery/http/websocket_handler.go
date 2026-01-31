package http

import (
	"net/http"

	"tsctl/internal/delivery/http/websocket"
	"tsctl/pkg/logger"

	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Log.WithError(err).Error("failed to upgrade websocket")
		return
	}

	client := websocket.NewClient(h.wsHub, conn)
	h.wsHub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
