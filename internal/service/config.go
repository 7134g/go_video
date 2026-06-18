package service

import (
	"errors"
	"fmt"
	"go_video/internal/controller"
	"go_video/internal/model"
	"go_video/internal/repository"
	"go_video/pkg/proxy"
	"os"
	"sync"
)

var (
	configService *ConfigService
	configSvcOnce sync.Once
)

type ConfigService struct {
	repo         *repository.ConfigRepository
	proxyServer  *proxy.Server
	proxyRunning bool
	mu           sync.Mutex
}

func GetConfigService() *ConfigService {
	configSvcOnce.Do(func() {
		configService = &ConfigService{
			repo: repository.GetConfigRepository(),
		}
	})
	return configService
}

func (s *ConfigService) Init() {
	cfg := s.repo.Get()
	if cfg.InterceptorEnabled {
		vpnAddress := cfg.VpnAddress
		if !cfg.VpnStatus {
			vpnAddress = ""
		}
		go s.startProxyServer(cfg.AgentAddress, vpnAddress)
	}
}

// Shutdown 在进程退出前关闭 MITM 代理（如已开启）。
func (s *ConfigService) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.proxyServer != nil {
		_ = s.proxyServer.Close()
		s.proxyServer = nil
		s.proxyRunning = false
	}
}

func (s *ConfigService) GetConfig() *model.Config {
	return s.repo.Get()
}

// SetFfmpegPromptDeclined 持久化“用户拒绝下载 ffmpeg”的选择，启动时据此跳过追问。
func (s *ConfigService) SetFfmpegPromptDeclined(v bool) error {
	cfg := s.repo.Get()
	cfg.FfmpegPromptDeclined = v
	return s.repo.Save(cfg)
}

func (s *ConfigService) UpdateConfig(updates map[string]interface{}) (*model.Config, error) {
	cfg := s.repo.Get()

	for key, val := range updates {
		switch key {
		case "max_concurrent_tasks":
			if v, ok := val.(float64); ok {
				if int(v) < 1 {
					return nil, errors.New("max_concurrent_tasks must be positive")
				}
				cfg.MaxConcurrentTasks = int(v)
			}
		case "max_segment_workers":
			if v, ok := val.(float64); ok {
				if int(v) < 1 {
					return nil, errors.New("max_segment_workers must be positive")
				}
				cfg.MaxSegmentWorkers = int(v)
			}
		case "download_dir":
			if v, ok := val.(string); ok {
				if err := os.MkdirAll(v, 0755); err != nil {
					return nil, errors.New("invalid download_dir path")
				}
				cfg.DownloadDir = v
			}
		case "max_consecutive_errors":
			if v, ok := val.(float64); ok {
				if int(v) < 1 {
					return nil, errors.New("max_consecutive_errors must be positive")
				}
				cfg.MaxConsecutiveErrors = int(v)
			}
		case "default_headers":
			if v, ok := val.(map[string]interface{}); ok {
				headers := make(map[string]string)
				for hk, hv := range v {
					if str, ok := hv.(string); ok {
						headers[hk] = str
					}
				}
				cfg.DefaultHeaders = headers
			}
		case "interceptor_enabled":
			if v, ok := val.(bool); ok {
				cfg.InterceptorEnabled = v
			}
		case "agent_address":
			if v, ok := val.(string); ok {
				cfg.AgentAddress = v
			}
		case "vpn_address":
			if v, ok := val.(string); ok {
				cfg.VpnAddress = v
			}
		case "vpn_status":
			if v, ok := val.(bool); ok {
				cfg.VpnStatus = v
			}
		}
	}

	if err := s.repo.Save(cfg); err != nil {
		return nil, err
	}

	// Apply config to download controller
	controller.GetController().ApplyConfig(
		cfg.DownloadDir,
		cfg.MaxConcurrentTasks,
		cfg.MaxSegmentWorkers,
		cfg.MaxConsecutiveErrors,
		cfg.DefaultHeaders,
	)

	// Handle interceptor
	if err := s.handleInterceptor(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (s *ConfigService) handleInterceptor(cfg *model.Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cfg.InterceptorEnabled && !s.proxyRunning {
		vpnAddress := cfg.VpnAddress
		if !cfg.VpnStatus {
			vpnAddress = ""
		}
		go s.startProxyServer(cfg.AgentAddress, vpnAddress)
		s.proxyRunning = true
	} else if !cfg.InterceptorEnabled && s.proxyRunning {
		if s.proxyServer != nil {
			if err := s.proxyServer.Close(); err != nil {
				return err
			}
		}
		s.proxyRunning = false
		s.proxyServer = nil
	}

	return nil
}

func (s *ConfigService) startProxyServer(agentAddress, vpnAddress string) {
	fmt.Println("开启被动代理 -> " + agentAddress)

	srv, err := proxy.NewServer(vpnAddress)
	if err != nil {
		fmt.Printf("代理服务器启动失败: %v\n", err)
		s.mu.Lock()
		s.proxyRunning = false
		s.mu.Unlock()
		return
	}
	s.proxyServer = srv

	go s.doTask(srv)

	err = srv.Listen(agentAddress)
	if err != nil {
		fmt.Printf("代理监听失败: %v\n", err)
		s.mu.Lock()
		s.proxyRunning = false
		s.mu.Unlock()
	}
}

func (s *ConfigService) doTask(srv *proxy.Server) {
	repo := repository.NewTaskRepository()
	for {
		select {
		case t := <-srv.Tasks():
			existing, err := repo.GetByURL(t.URL)
			switch {
			case repository.IsNotFound(err):
				_ = repo.Create(&model.Task{
					Name:      t.Title,
					URL:       t.URL,
					Header:    t.Headers,
					Type:      t.Type,
					CreatedAt: t.CreateAt,
				})
			case err == nil:
				// 已存在：只在新捕获更晚、且不在跑的情况下刷新名称和 Header。
				if existing.Status != model.TaskStatusRunning && t.CreateAt.After(existing.UpdatedAt) {
					_ = repo.UpdateNameAndHeader(existing.ID, t.Title, t.Headers)
				}
			default:
				fmt.Printf("代理任务查重失败: %v\n", err)
			}

		case <-srv.Stop:
			return
		}
	}
}
