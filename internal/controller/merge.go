package controller

import (
	"fmt"
	"luka/pkg/m3u8"
	"os"
	"path/filepath"
	"runtime"
)

func (c *DownloadController) mergeM3u8(task *DTask) error {
	BroadcastMessage(task.ID, "开始合并分片..."+task.Name)
	videoDir := safeJoin(c.config.DownloadDir, task.Name)

	ffmpegPath, err := c.ffmpegPath()
	if err != nil {
		return err
	}

	if err := m3u8.MergeFilesFfmpeg(videoDir, ffmpegPath); err != nil {
		return err
	}

	_ = os.RemoveAll(videoDir)
	BroadcastMessage(task.ID, "合并完成"+task.Name)
	return nil
}

// ffmpegPath 返回程序目录下 ffmpeg 可执行文件的路径，不存在时返回错误。
func (c *DownloadController) ffmpegPath() (string, error) {
	ffmpegName := "ffmpeg"
	if runtime.GOOS == "windows" {
		ffmpegName = "ffmpeg.exe"
	}
	ffmpegPath := filepath.Join(c.pwd, ffmpegName)
	if _, err := os.Stat(ffmpegPath); err != nil {
		return "", fmt.Errorf("未找到 ffmpeg，请下载 ffmpeg 放到程序目录 %s 后重试", c.pwd)
	}
	return ffmpegPath, nil
}
