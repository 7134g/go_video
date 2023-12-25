package task_control

import (
	"context"
	"crypto/tls"
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/db"
	"dv/internel/serve/api/internal/util/model"
	"log"
	"net/http"
	"net/url"
	"sync"
)

var (
	tcConfig *taskControlConfig
	errModel *model.ErrorModel
)

type TaskControl struct {
	Name string // 任务名

	wg        sync.WaitGroup
	mux       sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	doneCount uint // 完成数

	running   bool          // 是否正在运行
	vacancy   chan struct{} // 并发控制
	printStop chan struct{} // 打印进度
}

type taskControlConfig struct {
	// http
	Client  *http.Client
	Headers http.Header

	// 执行设置
	config.TaskControlConfig
}

func NewTaskControl(concurrency uint) *TaskControl {
	core := &TaskControl{
		wg:      sync.WaitGroup{},
		mux:     sync.Mutex{},
		running: false,
		vacancy: make(chan struct{}, concurrency),
	}
	return core
}

func InitTaskConfig(c config.Config) {
	errModel = model.NewErrorModel(db.GetDB())
	tcConfig = &taskControlConfig{
		Client:            &http.Client{Transport: getHttpProxy(c.HttpConfig)},
		Headers:           getHeader(c.HttpConfig),
		TaskControlConfig: c.TaskControlConfig,
	}
}

func getHttpProxy(c config.HttpConfig) http.RoundTripper {
	if !c.ProxyStatus {
		return nil
	}

	httpProxy := func(proxy string) func(*http.Request) (*url.URL, error) {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			log.Fatalln(err)
		}
		return http.ProxyURL(proxyUrl)
	}

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}, // 使用环境变量的代理
		Proxy: httpProxy(c.Proxy),
	}

	return t
}

func getHeader(c config.HttpConfig) http.Header {
	header := http.Header{}
	for k, v := range c.Headers {
		header.Set(k, v)
	}

	return header
}
