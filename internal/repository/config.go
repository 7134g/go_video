package repository

import (
	"encoding/json"
	"os"
	"sync"

	"go_video/internal/model"
)

const configFile = "config.json"

var (
	configRepo *ConfigRepository
	configOnce sync.Once
)

type ConfigRepository struct {
	mu       sync.RWMutex
	config   *model.Config
	filePath string
}

func GetConfigRepository() *ConfigRepository {
	configOnce.Do(func() {
		configRepo = &ConfigRepository{
			filePath: configFile,
		}
		configRepo.load()
	})
	return configRepo
}

func (r *ConfigRepository) load() {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		r.config = model.DefaultConfig()
		return
	}

	var cfg model.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		r.config = model.DefaultConfig()
		return
	}
	r.config = &cfg
}

func (r *ConfigRepository) Get() *model.Config {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cfg := *r.config
	return &cfg
}

func (r *ConfigRepository) Save(cfg *model.Config) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return err
	}

	r.config = cfg
	return nil
}
