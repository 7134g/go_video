package main

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/handler"
	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/svc/ws"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/task_serve.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	hub := ws.NewHub()
	hub.Run()
	server.AddRoute(ws.ShowLogWsRoute(hub))

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.PrintRoutes()
	server.Start()
}
