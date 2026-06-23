package api

import (
	"encoding/json"
	"net/http"
	"net/url"

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
//   - 连接建立时发送一次"所有任务进度快照"作为初始状态
//   - progressCh 接收 controller.BroadcastProgress 推来的单条进度
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

	// 连接建立时发送一次全量快照
	initProgress := ctrl.GetAllProgress()
	data, _ := json.Marshal(initProgress)
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return
	}

	progressCh := make(chan controller.ProgressInfo, 10)
	controller.AddProgressListener(progressCh)
	defer controller.RemoveProgressListener(progressCh)

	msgCh := make(chan controller.Message, 10)
	controller.AddMessageListener(msgCh)
	defer controller.RemoveMessageListener(msgCh)

	for {
		select {
		case info := <-progressCh:
			data, _ := json.Marshal(info)
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
