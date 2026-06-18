package model

import "github.com/gin-gonic/gin"

type Config struct {
	MaxConcurrentTasks   int               `json:"max_concurrent_tasks"`   // 并发任务数
	MaxSegmentWorkers    int               `json:"max_segment_workers"`    // 并发分片数
	DownloadDir          string            `json:"download_dir"`           // 下载地址
	MaxConsecutiveErrors int               `json:"max_consecutive_errors"` // 连续错误数
	DefaultHeaders       map[string]string `json:"default_headers"`        // 默认请求头
	InterceptorEnabled   bool              `json:"interceptor_enabled"`    // 是否开启拦截器
	AgentAddress         string            `json:"agent_address"`          // 拦截器代理地址
	VpnAddress           string            `json:"vpn_address"`            // vpn地址
	VpnStatus            bool              `json:"vpn_status"`             // 是否使用配置的vpn地址下载
	GinMode              string            `json:"gin_mode"`
	FfmpegPromptDeclined bool              `json:"ffmpeg_prompt_declined"` // 用户已拒绝下载 ffmpeg，启动时不再追问
}

func DefaultConfig() *Config {
	return &Config{
		MaxConcurrentTasks:   3,
		MaxSegmentWorkers:    5,
		DownloadDir:          "./downloads",
		MaxConsecutiveErrors: 10,
		DefaultHeaders: map[string]string{
			"user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		},
		InterceptorEnabled: false,
		AgentAddress:       "127.0.0.1:9999",
		VpnAddress:         "127.0.0.1:7890",
		GinMode:            gin.ReleaseMode,
	}
}
