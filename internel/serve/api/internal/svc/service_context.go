package svc

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/db"
	"dv/internel/serve/api/internal/middleware"
	"dv/internel/serve/api/internal/svc/task_control"
	"dv/internel/serve/api/internal/util/model"
	"dv/internel/serve/api/internal/util/proxy"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config          config.Config
	AuthInterceptor rest.Middleware

	TaskModel *model.TaskModel

	TaskControl *task_control.TaskControl
}

func NewServiceContext(c config.Config) *ServiceContext {
	db.InitSqlite(c.DB)
	task_control.InitTask(c)

	taskModel := model.NewTaskModel(db.GetDB())
	threading.GoSafe(func() {
		proxy.SetTaskDb(taskModel)
		proxy.SetServeProxyAddress(c.Proxy, "", "")
		proxy.OpenCert()
		proxy.SetMartianAddress(c.WebProxy)
		if err := proxy.Martian(); err != nil {
			panic(err)
		}
	})
	// 处理 ProxyCatchUrl 和 ProxyCatchHtml 匹配
	threading.GoSafe(func() {
		proxy.MatchInformation()
	})

	return &ServiceContext{
		Config:          c,
		AuthInterceptor: middleware.NewAuthInterceptorMiddleware().Handle,
		TaskModel:       taskModel,
		TaskControl:     task_control.NewTaskControl(c.TaskControlConfig.Concurrency),
	}
}
