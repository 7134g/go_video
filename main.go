package main

import (
	"context"
	"embed"
	"errors"
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
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist
var webFS embed.FS

func main() {
	if err := repository.InitDB(); err != nil {
		log.Fatal("Failed to init database:", err)
	}
	InitCa()

	svr := service.GetConfigService()
	cfg := svr.GetConfig()
	ensureFfmpeg(svr)
	controller.GetController().ApplyConfig(
		cfg.DownloadDir,
		cfg.MaxConcurrentTasks,
		cfg.MaxSegmentWorkers,
		cfg.MaxConsecutiveErrors,
		cfg.DefaultHeaders,
	)
	svr.Init()

	mode := cfg.GinMode
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	r := gin.Default()
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
		tasks.POST("/redownload", h.Redownload)
		tasks.GET("/progress", api.ProgressWS)
	}

	configHandler := api.NewConfigHandler()
	r.GET("/api/config", configHandler.Get)
	r.PUT("/api/config", configHandler.Update)

	ffmpegHandler := api.NewFfmpegHandler()
	r.GET("/api/ffmpeg/status", ffmpegHandler.Status)
	r.POST("/api/ffmpeg/download", ffmpegHandler.Download)

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

	httpSrv := &http.Server{Addr: "127.0.0.1:8080", Handler: r}

	go func() {
		fmt.Println("web地址 http://localhost:8080")
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("http server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	controller.GetController().StopAll()
	svr.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Println("http shutdown:", err)
	}
}

func InitCa() {
	installed, err := proxy.CheckCertInstalled()
	if err != nil {
		log.Fatal(err)
	}
	if installed {
		return
	}
	switch runtime.GOOS {
	case "windows":
		log.Fatal("CA 未安装,请运行 install_cert.exe")
	case "darwin":
		log.Fatal("CA 未安装,请运行 ./install_cert(将提示 sudo 密码)")
	case "linux":
		log.Fatal("CA 未安装,请运行 ./install_cert(会提示 sudo,并安装到系统及 NSS 库)")
	default:
		log.Fatal("CA 未安装,且当前平台不支持自动安装")
	}
}
