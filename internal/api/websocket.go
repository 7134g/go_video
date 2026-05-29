package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_video/internal/controller"
)

// 仅允许 localhost 同源 WebSocket 升级，避免任意网页接管本地下载器。
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true // 非浏览器客户端 (curl/python) 无 Origin，放行
		}
		u, err := url.Parse(origin)
		if err != nil {
			return false
		}
		switch u.Hostname() {
		case "localhost", "127.0.0.1", "::1":
			return true
		}
		return false
	},
}

// ProgressWS 是 /api/tasks/progress 的 WebSocket 端点。
// 每个连接两路输出：
//   - ticker 每秒推一次"所有任务进度快照"，用于前端进度条
//   - msgCh 接收 controller.BroadcastMessage 推来的事件文案
//
// 任一写失败即视为客户端断开，整个 handler 返回，监听器随 defer 注销。
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
