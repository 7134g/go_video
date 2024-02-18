package ws_conn

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// DefaultServeWs ServeWs handles websocket requests from the peer.
func DefaultServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) error {
	client, err := InitClient(hub, w, r)
	if err != nil {
		return err
	}

	Run(client)
	return nil
}

func InitClient(hub *Hub, w http.ResponseWriter, r *http.Request) (*Client, error) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	client := NewClient(hub, conn)

	return client, nil
}

func Run(client *Client) {
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
