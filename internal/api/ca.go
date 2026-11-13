package api

import (
	"net/http"

	"luka/pkg/proxy"

	"github.com/gin-gonic/gin"
)

type CaHandler struct{}

func NewCaHandler() *CaHandler { return &CaHandler{} }

func (h *CaHandler) Status(c *gin.Context) {
	installed, err := proxy.CheckCertInstalled()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"installed": installed})
}
