package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	cfg ProjectConfig
)

func GetConfig() *ProjectConfig {
	return &cfg
}

type ProjectConfig struct {
	SaveDir    string `yaml:"save_dir"`
	HistoryDir string `yaml:"history_dir"`
	LogDir     string `yaml:"log_dir"`
	LogStatus  bool   `yaml:"log_status"`

	TaskList          string `yaml:"task_list"`            // 任务清单
	TaskErrorMaxCount uint   `yaml:"task_error_max_count"` // 任务最大数
	TaskErrorDuration uint   `yaml:"task_error_duration"`  // 错误时候休眠多久后重试
	Concurrency       uint   `yaml:"concurrency"`          // 并发数
	ConcurrencyM3u8   uint   `yaml:"concurrency_m3u8"`     // m3u8 片段并发大小
	TaskClear         bool   `yaml:"task_clear"`           // 清空任务清单文件

	Headers     map[string]string `yaml:"headers"` // 请求头
	Proxy       string            `yaml:"proxy"`
	ProxyStatus bool              `yaml:"proxy_status"` // 代理是否开启的状态
}

func LoadConfig() {
	// 读取 YAML 文件内容
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 解析 YAML
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	configInit()
}

// 初始化
func configInit() {
	Header = createHeader()

	if _, err := os.Stat(cfg.HistoryDir); err != nil {
		_ = os.MkdirAll(cfg.HistoryDir, os.ModeDir)
	}
	if _, err := os.Stat(cfg.SaveDir); err != nil {
		cfg.SaveDir = "./download"
		_ = os.MkdirAll(cfg.SaveDir, os.ModeDir)
	}

	if cfg.ProxyStatus {
		Client = GetHttpProxyClient(cfg.Proxy)
	}

	if cfg.LogStatus {
		if _, err := os.Stat(cfg.LogDir); err != nil {
			_ = os.MkdirAll(cfg.LogDir, os.ModeDir)
		}
		nowTime := time.Now().Format("2006_01_02_15_04_05")
		logName := filepath.Join(cfg.LogDir, nowTime+".log")
		f, _ := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		log.SetOutput(io.MultiWriter(os.Stdout, f))
	}
	log.SetFlags(log.Ltime)
	log.Printf("读取 %s 文件, 同时进行任务最大值为 %d , 操作目录为 %s \n",
		cfg.TaskList, cfg.Concurrency, cfg.SaveDir)
}

var (
	Header http.Header
)

func createHeader() http.Header {
	headers := make(http.Header, 0)
	for k, v := range cfg.Headers {
		headers.Set(k, v)
	}
	return headers
}

func createHeaderM3U8() string {
	var headers string
	for k, v := range cfg.Headers {
		if len(headers) != 0 {
			headers += "|"
		}
		headers += fmt.Sprintf("%v:%v", k, v)
	}
	return headers
}
