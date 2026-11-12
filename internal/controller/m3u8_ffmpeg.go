package controller

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// downloadM3u8WithFfmpeg 让 ffmpeg 直接拉取 m3u8 并封装为 <download_dir>/<任务名>.mp4。
// 不做分片级断点续传与重试，失败原因（含 stderr）记录到日志。
func (c *DownloadController) downloadM3u8WithFfmpeg(task *DTask) error {
	BroadcastMessage(task.ID, "ffmpeg 开始下载..."+task.Name)

	ffmpegPath, err := c.ffmpegPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(c.config.DownloadDir, 0755); err != nil {
		return err
	}

	task.Progress.SetTotal(1)

	outPath := safeJoin(c.config.DownloadDir, task.Name) + ".mp4"
	args := []string{"-y", "-nostdin"}
	if headers := ffmpegHeaders(task.Header, c.config.DefaultHeaders); headers != "" {
		args = append(args, "-headers", headers)
	}
	if proxyURL := c.ffmpegProxyURL(); proxyURL != "" {
		args = append(args, "-http_proxy", proxyURL)
	}
	args = append(args, "-i", task.URL, "-c", "copy", outPath)

	cmd := exec.CommandContext(task.ctx, ffmpegPath, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if task.ctx.Err() != nil {
		return context.Canceled
	}
	if err != nil {
		slog.Error("ffmpeg 直连下载失败", "task", task.Name, "url", task.URL, "error", err, "stderr", strings.TrimSpace(stderr.String()))
		return fmt.Errorf("ffmpeg 下载失败: %w", err)
	}

	task.Progress.AddDone(1)
	BroadcastMessage(task.ID, "下载完成"+task.Name)
	return nil
}

// ffmpegHeaders 把 Header 拼成 ffmpeg -headers 参数要求的格式（每行以 \r\n 结尾）。
// 任务 Header 为空时使用默认 Header。
func ffmpegHeaders(taskHeader http.Header, defaultHeaders map[string]string) string {
	source := taskHeader
	if len(source) == 0 {
		source = convertHeaders(defaultHeaders)
	}
	var b strings.Builder
	for k, values := range source {
		for _, v := range values {
			b.WriteString(k + ": " + v + "\r\n")
		}
	}
	return b.String()
}

// ffmpegProxyURL 返回 ffmpeg -http_proxy 使用的代理地址，未启用时返回空字符串。
func (c *DownloadController) ffmpegProxyURL() string {
	if !c.config.VpnStatus || c.config.VpnAddress == "" {
		return ""
	}
	return "http://" + c.config.VpnAddress
}
