package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_video/internal/controller"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ProgressWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctrl := controller.GetController()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	msgCh := make(chan controller.Message, 10)
	controller.AddMessageListener(msgCh)
	defer controller.RemoveMessageListener(msgCh)

	for {
		select {
		case <-ticker.C:
			progress := ctrl.GetAllProgress()
			data, _ := json.Marshal(progress)
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case msg := <-msgCh:
			data, _ := json.Marshal(msg)
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		}
	}
}
