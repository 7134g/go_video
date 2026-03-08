package main

import (
	"embed"
	"go_video/internal/api"
	"go_video/internal/controller"
	"go_video/internal/repository"
	"go_video/internal/service"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist
var webFS embed.FS

func main() {
	if err := repository.InitDB(); err != nil {
		log.Fatal("Failed to init database:", err)
	}

	// 加载配置并应用到 controller
	cfg := service.GetConfigService().GetConfig()
	controller.GetController().ApplyConfig(
		cfg.DownloadDir,
		cfg.MaxConcurrentTasks,
		cfg.MaxSegmentWorkers,
		cfg.MaxConsecutiveErrors,
		cfg.DefaultHeaders,
	)

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	h := api.NewTaskHandler()

	tasks := r.Group("/api/tasks")
	{
		tasks.GET("", h.List)
		tasks.POST("", h.Create)
		tasks.DELETE("/:id", h.Delete)
		tasks.PUT("/:id", h.Update)
		tasks.POST("/start", h.Start)
		tasks.POST("/:id/pause", h.Pause)
		tasks.POST("/:id/retry", h.Retry)
		tasks.GET("/progress", api.ProgressWS)
	}

	configHandler := api.NewConfigHandler()
	r.GET("/api/config", configHandler.Get)
	r.PUT("/api/config", configHandler.Update)

	distFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Fatal("Failed to load web files:", err)
	}
	r.NoRoute(func(c *gin.Context) {
		file, err := distFS.Open(c.Request.URL.Path[1:])
		if err != nil {
			data, _ := fs.ReadFile(distFS, "index.html")
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			return
		}
		_ = file.Close()
		c.FileFromFS(c.Request.URL.Path, http.FS(distFS))
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		os.Exit(0)
	}()

	_ = r.Run(":8080")
}
