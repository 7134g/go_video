//go:build !desktop && !bindings

package main

import (
	"context"
	"embed"
	"errors"
	"go_video/internal/api"
	"go_video/internal/controller"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist
var webFS embed.FS

func main() {
	svr := initShared()

	cfg := svr.GetConfig()
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

	caHandler := api.NewCaHandler()
	r.GET("/api/ca/status", caHandler.Status)

	distFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		slog.Error("加载 web 文件失败", "error", err)
		os.Exit(1)
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
		slog.Info("web地址 http://localhost:8080")
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down...")
	controller.GetController().StopAll()
	svr.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		slog.Warn("http shutdown error", "error", err)
	}
}
