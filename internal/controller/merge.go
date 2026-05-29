package controller

import (
	"fmt"
	"go_video/pkg/m3u8"
	"os"
	"path/filepath"
	"runtime"
)

func (c *DownloadController) mergeM3u8(task *DTask) error {
	BroadcastMessage(task.ID, "开始合并分片..."+task.Name)
	videoDir := safeJoin(c.config.DownloadDir, task.Name)

	// 优先纯 Go remux，无需 ffmpeg。
	nativeErr := m3u8.MergeFilesNative(videoDir)
	if nativeErr == nil {
		_ = os.RemoveAll(videoDir)
		BroadcastMessage(task.ID, "合并完成"+task.Name)
		return nil
	}

	// 回退到 ffmpeg。
	var ffmpegName string
	if runtime.GOOS == "windows" {
		ffmpegName = "ffmpeg.exe"
	} else {
		ffmpegName = "ffmpeg"
	}
	ffmpegPath := filepath.Join(c.pwd, ffmpegName)
	if _, err := os.Stat(ffmpegPath); err != nil {
		return fmt.Errorf("原生合并失败(%v)且未找到 ffmpeg，请下载 ffmpeg 放到程序目录 %s 后重试", nativeErr, c.pwd)
	}

	BroadcastMessage(task.ID, "原生合并失败，尝试 ffmpeg..."+nativeErr.Error())
	if err := m3u8.MergeFilesFfmpeg(videoDir, ffmpegPath); err != nil {
		return err
	}

	_ = os.RemoveAll(videoDir)
	BroadcastMessage(task.ID, "合并完成"+task.Name)
	return nil
}
