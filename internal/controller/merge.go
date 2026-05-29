package controller

import (
	"go_video/pkg/m3u8"
	"os"
	"path/filepath"
	"runtime"
)

func (c *DownloadController) mergeM3u8(task *DTask) error {
	BroadcastMessage(task.ID, "开始合并分片..."+task.Name)
	var ffmpegName string
	if runtime.GOOS == "windows" {
		ffmpegName = "ffmpeg.exe"
	} else {
		ffmpegName = "ffmpeg"
	}
	ffmpegPath := filepath.Join(c.pwd, ffmpegName)
	videoDir := safeJoin(c.config.DownloadDir, task.Name)
	err := m3u8.MergeFilesFfmpeg(videoDir, ffmpegPath)
	if err != nil {
		return err
	}

	_ = os.RemoveAll(videoDir)

	BroadcastMessage(task.ID, "合并完成"+task.Name)
	return nil
}
