package handler

import (
	"dv/internel/serve/api/internal/handler/ws"
	"dv/internel/serve/api/internal/svc"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

func RegisterWSHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.AuthInterceptor},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/log",
					Handler: ws.LogHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/ws"),
	)
}
