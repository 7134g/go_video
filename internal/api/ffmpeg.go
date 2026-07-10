package api

import (
	"go_video/internal/downloader"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FfmpegHandler struct{}

func NewFfmpegHandler() *FfmpegHandler { return &FfmpegHandler{} }

// Status 返回 ffmpeg 是否已存在、当前平台是否支持自动下载。
func (h *FfmpegHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"exists":    downloader.Exists(),
		"supported": downloader.Supported(),
	})
}

// Download 同步下载 ffmpeg 到程序目录（耗时较长，前端需放宽超时）。
func (h *FfmpegHandler) Download(c *gin.Context) {
	if err := downloader.Download(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"exists": true})
}
