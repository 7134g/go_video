package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	DB string // 数据库

	HttpConfig

	TaskControlConfig
}

type HttpConfig struct {
	Headers     map[string]string // 默认请求头
	Proxy       string            // 代理地址
	ProxyStatus bool              // 代理是否开启
}

type TaskControlConfig struct {
	Concurrency       uint   // 并发数
	ConcurrencyM3u8   uint   // m3u8 片段并发大小
	SaveDir           string // 存储位置
	TaskErrorMaxCount uint   // 任务连续最大错误次数
	TaskErrorDuration uint   // 错误时候休眠多久后重试(秒)
}
