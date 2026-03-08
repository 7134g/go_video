package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go_video/internal/service"
)

type ConfigHandler struct {
	svc *service.ConfigService
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{svc: service.GetConfigService()}
}

func (h *ConfigHandler) Get(c *gin.Context) {
	cfg := h.svc.GetConfig()
	c.JSON(http.StatusOK, cfg)
}

func (h *ConfigHandler) Update(c *gin.Context) {
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg, err := h.svc.UpdateConfig(updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}
