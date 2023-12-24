package svc

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/db"
	"dv/internel/serve/api/internal/middleware"
	"dv/internel/serve/api/internal/svc/task_control"
	model2 "dv/internel/serve/api/internal/util/model"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config          config.Config
	AuthInterceptor rest.Middleware

	TaskModel  *model2.TaskModel
	ErrorModel *model2.ErrorModel

	TaskControl *task_control.TaskControl
}

func NewServiceContext(c config.Config) *ServiceContext {
	db.InitSqlite(c.DB)

	return &ServiceContext{
		Config:          c,
		AuthInterceptor: middleware.NewAuthInterceptorMiddleware().Handle,
		TaskModel:       model2.NewTaskModel(db.GetDB()),
		ErrorModel:      model2.NewErrorModel(db.GetDB()),
		TaskControl:     task_control.NewTaskControl(c.TaskControlConfig.Concurrency),
	}
}