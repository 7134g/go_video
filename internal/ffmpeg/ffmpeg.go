package ffmpeg

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Exists() bool {
	_, err := os.Stat(Name())
	return err == nil
}

func Name() string {
	if runtime.GOOS == "windows" {
		return "ffmpeg.exe"
	}
	return "ffmpeg"
}

func Supported() bool {
	// Windows x86_64 and arm64, macOS x86_64 and arm64
	switch runtime.GOOS + "/" + runtime.GOARCH {
	case "windows/amd64", "windows/arm64", "darwin/amd64", "darwin/arm64":
		return true
	default:
		return false
	}
}

func URL() string {
	switch {
	case runtime.GOOS == "windows":
		// BtbN FFmpeg builds — Windows only ships as zip
		return "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip"
	case runtime.GOOS == "darwin":
		// evermeet.cx — static ffmpeg binary for macOS
		return "https://evermeet.cx/ffmpeg/ffmpeg.zip"
	default:
		return ""
	}
}

func Download(ctx context.Context) error {
	url := URL()
	if url == "" {
		return fmt.Errorf("当前平台 %s/%s 不支持自动下载 ffmpeg", runtime.GOOS, runtime.GOARCH)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败: HTTP %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "ffmpeg-*.zip")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return err
	}
	tmpFile.Close()

	reader, err := zip.OpenReader(tmpPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		base := filepath.Base(f.Name)
		if !strings.EqualFold(base, "ffmpeg.exe") && !strings.EqualFold(base, "ffmpeg") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		out, err := os.OpenFile(Name(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, rc); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("在压缩包中未找到 ffmpeg 可执行文件")
}
