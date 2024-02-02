package ws

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

func ShowLogWsRoute(hub *Hub) rest.Route {
	return rest.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			client, err := InitClient(hub, w, r)
			if err != nil {
				logx.Error(err)
				return
			}
			client.SetReadHandle(func(message []byte) []byte {
				var b []byte
				var err error
				switch string(message) {
				case "get":
					b, err = json.Marshal(messageResponse{
						Code: 0,
						Msg:  "",
						Data: message, // todo 日志信息
					})
					if err != nil {
						logx.Error(err)
						return nil
					}
				}

				return b
			})

			//client.SetWriteHandle(func(message []byte) []byte {
			//	return message
			//})
			Run(client)

		},
	}
}
