package main

import (
	"bufio"
	"encoding/json"
	"go_video/internal/logger"
	"go_video/pkg/proxy"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strings"

	"go_video/internal/controller"
	"go_video/internal/model"
	"go_video/internal/repository"
	"go_video/internal/service"
)

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

// initShared 执行 web 和 desktop 模式共享的初始化逻辑。
func initShared() *service.ConfigService {
	if err := repository.InitDB(); err != nil {
		log.Fatal("Failed to init database:", err)
	}

	svr := service.GetConfigService()
	cfg := svr.GetConfig()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}
	logger.InitWithLevel(cfg.LogLevel, pwd+"/"+logger.LogDir)

	controller.GetController().ApplyConfig(
		cfg.DownloadDir,
		cfg.MaxConcurrentTasks,
		cfg.MaxSegmentWorkers,
		cfg.MaxConsecutiveErrors,
		cfg.DefaultHeaders,
	)
	importTaskFile(cfg.DefaultHeaders)
	svr.Init()
	return svr
}

func importTaskFile(defaultHeaders map[string]string) {
	f, err := os.Open("task.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		slog.Warn("读取 task.txt 失败", "error", err)
		return
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	f.Close()

	repo := repository.NewTaskRepository()

	for i := 0; i+1 < len(lines); i += 2 {
		name := lines[i]
		url := lines[i+1]

		if _, err := repo.GetByURL(url); err == nil {
			slog.Info("task URL 已存在，跳过", "name", name)
			continue
		}

		taskType := "m3u8"
		if strings.HasSuffix(url, ".mp4") {
			taskType = "mp4"
		}

		h := make(http.Header)
		for k, v := range defaultHeaders {
			h.Set(k, v)
		}
		headerJSON, _ := json.Marshal(h)

		task := &model.Task{
			Name:   name,
			URL:    url,
			Header: string(headerJSON),
			Type:   taskType,
			Status: model.TaskStatusPending,
		}

		if err := repo.Create(task); err != nil {
			slog.Error("导入任务失败", "name", name, "error", err)
		}
	}

	if err := os.Remove("task.txt"); err != nil {
		slog.Warn("删除 task.txt 失败", "error", err)
	}
}
