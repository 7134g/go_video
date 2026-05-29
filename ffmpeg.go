package main

import (
	"bufio"
	"context"
	"fmt"
	"go_video/internal/ffmpeg"
	"go_video/internal/service"
	"os"
	"strings"
)

// ensureFfmpeg 在启动时检查程序目录下是否存在 ffmpeg；不存在则询问用户是否下载。
// 用户选择“否”会被记住（写入 config.json），下次启动不再追问；成功下载后文件已存在，
// 自然不会再问。ffmpeg 仅作合并兜底：默认用纯 Go remux，遇到非 H264/H265 + AAC/MP3 的片源才需要它。
func ensureFfmpeg(svr *service.ConfigService) {
	if ffmpeg.Exists() {
		return
	}
	if svr.GetConfig().FfmpegPromptDeclined {
		return // 用户此前已选择不下载，不再追问
	}

	name := ffmpeg.Name()
	if !ffmpeg.Supported() {
		fmt.Printf("未检测到 %s，且当前平台不支持自动下载，部分视频格式可能无法合并。\n", name)
		return
	}

	fmt.Printf("未检测到 %s。部分视频格式（非 H264/H265 + AAC/MP3）需要 ffmpeg 才能合并。\n", name)
	fmt.Print("是否现在下载 ffmpeg 到当前目录？[y/N]: ")
	// 非交互启动（无 tty）时 ReadString 立即返回 EOF，按 “否” 处理，不阻塞。
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if ans := strings.ToLower(strings.TrimSpace(line)); ans != "y" && ans != "yes" {
		_ = svr.SetFfmpegPromptDeclined(true)
		fmt.Println("已记住选择：不下载 ffmpeg。如需启用，可在 Web 配置页点击下载，或删除 config.json 中的 ffmpeg_prompt_declined 后重启。")
		return
	}

	fmt.Printf("正在下载 ffmpeg: %s\n", ffmpeg.URL())
	if err := ffmpeg.Download(context.Background()); err != nil {
		fmt.Printf("下载失败: %v\n请手动从 %s 下载 %s 放到程序目录\n", err, ffmpeg.URL(), name)
		return
	}
	fmt.Printf("ffmpeg 已安装。\n")
}
