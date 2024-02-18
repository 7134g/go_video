package ws_conn

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
	"testing"
)

func TestZeroWs(t *testing.T) {
	h := NewHub()
	go h.Run()

	engine := rest.MustNewServer(rest.RestConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				Mode: "console",
			},
		},
		Host:         "localhost",
		Port:         10999,
		Timeout:      10000,
		CpuThreshold: 500,
	})
	defer engine.Stop()

	engine.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		},
	})

	engine.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/log",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			client, err := InitClient(h, w, r)
			if err != nil {
				logx.Error(err)
				return
			}
			client.SetReadHandle(func(message []byte) []byte {
				return append([]byte("receive a message "), message...)
			})

			client.SetWriteHandle(func(message []byte) []byte {
				return append(message, []byte(", write back")...)
			})
			Run(client)
		},
	})

	engine.PrintRoutes()
	engine.Start()
}
