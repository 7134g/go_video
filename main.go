package main

import (
	"embed"
	"fmt"
	"go_video/internal/api"
	"go_video/internal/controller"
	"go_video/internal/repository"
	"go_video/internal/service"
	"go_video/pkg/proxy"
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
	InitCa()

	// 加载配置并应用到 controller
	svr := service.GetConfigService()
	cfg := svr.GetConfig()
	controller.GetController().ApplyConfig(
		cfg.DownloadDir,
		cfg.MaxConcurrentTasks,
		cfg.MaxSegmentWorkers,
		cfg.MaxConsecutiveErrors,
		cfg.DefaultHeaders,
	)
	svr.Init()

	fmt.Println("web地址 http://localhost:8080")
	r := gin.Default()
	if cfg.GinMode == "" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	h := api.NewTaskHandler()

	tasks := r.Group("/api/tasks")
	{
		tasks.GET("", h.List)
		tasks.POST("", h.Create)
		tasks.POST("/delete", h.Delete)
		tasks.POST("/update", h.Update)
		tasks.POST("/start", h.Start)
		tasks.POST("/pause", h.Pause)
		tasks.POST("/retry", h.Retry)
		tasks.POST("/start-one", h.StartOne)
		tasks.POST("/stop-all", h.PauseAll)
		tasks.POST("/update-title", h.UpdateTitle)
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

func InitCa() {
	installed, err := proxy.CheckCertInstalled()
	if err != nil {
		panic(err)
	}
	if !installed {
		panic("需要先安装证书")
	}
}
