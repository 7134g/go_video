package svc

import (
	"bytes"
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/db"
	"dv/internel/serve/api/internal/middleware"
	"dv/internel/serve/api/internal/svc/task_control"
	"dv/internel/serve/api/internal/util/files"
	"dv/internel/serve/api/internal/util/model"
	"dv/internel/serve/api/internal/util/proxy"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/rest"
	"io"
	"os"
	"time"
)

type ServiceContext struct {
	Config          config.Config
	AuthInterceptor rest.Middleware
	TaskModel       *model.TaskModel
	TaskControl     *task_control.TaskControl
	LogData         *bytes.Buffer
}

func NewServiceContext(c config.Config) *ServiceContext {
	taskModel := model.NewTaskModel(db.InitSqlite(c.DB))

	task_control.InitTask(c)

	// 开启被动代理
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

	// 设置日志
	f, err := files.GetFile(fmt.Sprintf("./log/%s.log", time.Now().Format(time.DateOnly)))
	if err != nil {
		panic(err)
	}
	logData := bytes.NewBuffer(nil)
	logWrite := logx.NewWriter(io.MultiWriter(os.Stdout, logData, f))
	logx.SetWriter(logWrite)

	return &ServiceContext{
		Config:          c,
		AuthInterceptor: middleware.NewAuthInterceptorMiddleware().Handle,
		TaskModel:       taskModel,
		TaskControl:     task_control.NewTaskControl(c.TaskControlConfig.Concurrency),
		LogData:         logData,
	}
}
