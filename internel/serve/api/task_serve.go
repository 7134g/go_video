package main

import (
	"dv/internel/serve/api/internal/svc/task_control"
	"flag"
	"fmt"

	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/handler"
	"dv/internel/serve/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/task_serve.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	task_control.InitTaskConfig(c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.PrintRoutes()
	server.Start()
}
