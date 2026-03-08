package service

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"go_video/internal/controller"
	"go_video/internal/model"
	"go_video/internal/repository"
	"go_video/pkg/proxy"
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
		go s.startProxyServer(cfg.ProxyAddress)
	}
}

func (s *ConfigService) GetConfig() *model.Config {
	return s.repo.Get()
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
		case "proxy_address":
			if v, ok := val.(string); ok {
				cfg.ProxyAddress = v
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
		go s.startProxyServer(cfg.ProxyAddress)
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

func (s *ConfigService) startProxyServer(address string) {
	fmt.Println("开启代理" + address)

	srv, err := proxy.NewServer()
	if err != nil {
		panic(err)
	}
	s.proxyServer = srv
	err = srv.Listen(address)
	if err != nil {
		panic(err)
	}
}
